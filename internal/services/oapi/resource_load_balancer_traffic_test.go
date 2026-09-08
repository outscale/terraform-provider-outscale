package oapi_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	lbTrafficBackendPort   = 8000
	lbTrafficClientPort    = 80
	lbTrafficRequestCount  = 30
	lbTrafficBackendCount  = 2
	lbTrafficHealthTimeout = 10 * time.Minute
)

func TestAccVM_LoadBalancer_HTTPBackends(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	lbResourceName := "outscale_load_balancer.lbu_traffic"
	// Load balancer names are limited to 32 chars; RandomWithPrefix adds "-<int64>".
	lbName := acctest.RandomWithPrefix("testacc-lbt")
	sgName := acctest.RandomWithPrefix("testacc-sg")
	region := utils.GetRegion()
	vmType := testAccVmType

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckLoadBalancerTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerTrafficConfig(lbName, omi, region, vmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerHTTPBackends(t, lbResourceName, lbTrafficBackendCount),
				),
			},
		},
	})
}

func testAccLoadBalancerTrafficConfig(lbName, omi, region, vmType, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_security_group" "sg_lb_traffic" {
  security_group_name = "%[5]s"
  description         = "Used in the terraform acceptance tests"
}

resource "outscale_security_group_rule" "lb_traffic_http" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.sg_lb_traffic.security_group_id
  from_port_range   = %[6]d
  to_port_range     = %[6]d
  ip_protocol       = "tcp"
  ip_range          = "0.0.0.0/0"
}

resource "outscale_vm" "lb_traffic_backend" {
  count              = %[7]d
  image_id           = "%[3]s"
  vm_type            = "%[4]s"
  security_group_ids = [outscale_security_group.sg_lb_traffic.security_group_id]
  user_data = base64encode(<<-EOT
#!/bin/bash
set -e
mkdir -p /var/www/html
echo "backend-${count.index}" > /var/www/html/index.html
if command -v python3 >/dev/null 2>&1; then
  nohup python3 -m http.server %[6]d --bind 0.0.0.0 --directory /var/www/html >/var/log/lb-traffic-http.log 2>&1 &
else
  cd /var/www/html
  nohup python -m SimpleHTTPServer %[6]d >/var/log/lb-traffic-http.log 2>&1 &
fi
EOT
  )
}

resource "outscale_load_balancer" "lbu_traffic" {
  load_balancer_name = "%[1]s"
  subregion_names    = ["%[2]sa"]
  listeners {
    backend_port           = %[6]d
    backend_protocol       = "HTTP"
    load_balancer_port     = %[8]d
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_attributes" "lbu_traffic" {
  load_balancer_name = outscale_load_balancer.lbu_traffic.load_balancer_name
  health_check {
    healthy_threshold   = 2
    check_interval      = 10
    path                = "/"
    port                = %[6]d
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

resource "outscale_load_balancer_vms" "lbu_traffic_backends" {
  load_balancer_name = outscale_load_balancer.lbu_traffic.load_balancer_name
  backend_vm_ids     = outscale_vm.lb_traffic_backend[*].vm_id
  depends_on = [
    outscale_load_balancer_attributes.lbu_traffic,
    outscale_security_group_rule.lb_traffic_http,
  ]
}
`, lbName, region, omi, vmType, sgName, lbTrafficBackendPort, lbTrafficBackendCount, lbTrafficClientPort)
}

func testAccCheckLoadBalancerTrafficDestroy(s *terraform.State) error {
	if testacc.ConfiguredClient == nil || testacc.ConfiguredClient.OSCAPI == nil {
		return fmt.Errorf("configured API client is nil")
	}
	conn := testacc.ConfiguredClient.OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		var resp oscgo.ReadLoadBalancersResponse
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			req := &oscgo.ReadLoadBalancersRequest{
				Filters: &oscgo.FiltersLoadBalancer{
					LoadBalancerNames: &[]string{rs.Primary.ID},
				},
			}

			rp, httpResp, err := conn.LoadBalancerApi.ReadLoadBalancers(context.Background()).ReadLoadBalancersRequest(*req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err == nil {
			if len(*resp.LoadBalancers) != 0 &&
				*(*resp.LoadBalancers)[0].LoadBalancerName == rs.Primary.ID {
				return fmt.Errorf("load balancer still exists: %s", rs.Primary.ID)
			}
		}

		if strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound") {
			return nil
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckLoadBalancerHTTPBackends(t *testing.T, lbResourceName string, backendCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[lbResourceName]
		if !ok {
			return fmt.Errorf("not found: %s", lbResourceName)
		}

		dnsName := rs.Primary.Attributes["dns_name"]
		lbName := rs.Primary.Attributes["load_balancer_name"]
		if dnsName == "" || lbName == "" {
			return fmt.Errorf("load balancer dns_name or load_balancer_name is not set")
		}

		conn := testacc.ConfiguredClient.OSCAPI
		if conn == nil {
			return fmt.Errorf("configured API client is nil")
		}

		if err := waitForLoadBalancerBackendsHealthy(t, conn, lbName, backendCount); err != nil {
			return err
		}

		return verifyLoadBalancerHTTPDistribution(dnsName, backendCount)
	}
}

func waitForLoadBalancerBackendsHealthy(t *testing.T, conn *oscgo.APIClient, lbName string, backendCount int) error {
	req := oscgo.ReadVmsHealthRequest{
		LoadBalancerName: lbName,
	}

	attempt := 0
	start := time.Now()
	var lastLogElapsed time.Duration

	return retry.RetryContext(context.Background(), lbTrafficHealthTimeout, func() *retry.RetryError {
		attempt++
		elapsed := time.Since(start)

		resp, httpResp, err := conn.LoadBalancerApi.ReadVmsHealth(context.Background()).ReadVmsHealthRequest(req).Execute()
		if err != nil {
			if attempt == 1 || elapsed-lastLogElapsed >= 30*time.Second {
				t.Logf("waiting for healthy backends on %s: attempt %d (elapsed %s): API error: %v",
					lbName, attempt, elapsed.Round(time.Second), err)
				lastLogElapsed = elapsed
			}
			return utils.CheckThrottling(httpResp, err)
		}

		health := resp.GetBackendVmHealth()
		healthyCount := 0
		for _, backend := range health {
			if isLoadBalancerBackendHealthy(backend.GetState()) {
				healthyCount++
			}
		}

		if len(health) < backendCount || healthyCount < backendCount {
			if attempt == 1 || elapsed-lastLogElapsed >= 30*time.Second {
				t.Logf("waiting for healthy backends on %s: attempt %d (elapsed %s): %d/%d healthy, %d/%d entries",
					lbName, attempt, elapsed.Round(time.Second), healthyCount, backendCount, len(health), backendCount)
				lastLogElapsed = elapsed
			}
			if len(health) < backendCount {
				return retry.RetryableError(fmt.Errorf("expected %d backend health entries, got %d", backendCount, len(health)))
			}
			return retry.RetryableError(fmt.Errorf("expected %d healthy backends, got %d", backendCount, healthyCount))
		}

		return nil
	})
}

func isLoadBalancerBackendHealthy(state string) bool {
	switch strings.ToUpper(state) {
	case "UP", "INSERVICE":
		return true
	default:
		return false
	}
}

func verifyLoadBalancerHTTPDistribution(dnsName string, backendCount int) error {
	client := &http.Client{Timeout: 30 * time.Second}
	url := fmt.Sprintf("http://%s:%d/", dnsName, lbTrafficClientPort)

	seen := make(map[string]struct{})
	for i := 0; i < lbTrafficRequestCount; i++ {
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("http request %d failed: %w", i+1, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("read http response %d: %w", i+1, err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http request %d: expected status 200, got %d (body: %q)", i+1, resp.StatusCode, strings.TrimSpace(string(body)))
		}

		backendID := strings.TrimSpace(string(body))
		if backendID == "" {
			return fmt.Errorf("http request %d: empty response body", i+1)
		}

		seen[backendID] = struct{}{}
	}

	for i := 0; i < backendCount; i++ {
		expected := fmt.Sprintf("backend-%d", i)
		if _, ok := seen[expected]; !ok {
			seenBackends := make([]string, 0, len(seen))
			for backend := range seen {
				seenBackends = append(seenBackends, backend)
			}
			return fmt.Errorf("did not receive traffic from %q after %d requests; saw responses: %v", expected, lbTrafficRequestCount, seenBackends)
		}
	}

	return nil
}
