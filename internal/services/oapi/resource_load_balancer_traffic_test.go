package oapi_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

const (
	lbTrafficHTTPPort     = 80
	lbTrafficBackendCount = 3
	lbTrafficRequestCount = 30

	lbTrafficHealthTimeout = 10 * time.Minute
	lbTrafficHTTPRetry     = 2 * time.Minute
)

var lbTrafficHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
	},
}

func TestAccVM_LoadBalancer_HTTPBackends(t *testing.T) {
	lbName := acctest.RandomWithPrefix("testacc-lbt")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	testacc.Test(t, resource.TestCase{
		CheckDestroy: testAccCheckLoadBalancerTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerTrafficConfig(lbName, sgName),
				Check:  testAccCheckLoadBalancerTraffic(t),
			},
		},
	})
}

func testAccCheckLoadBalancerTraffic(t *testing.T) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		lbRS := s.RootModule().Resources["outscale_load_balancer.lbu_traffic"]
		if lbRS == nil {
			return fmt.Errorf("not found: outscale_load_balancer.lbu_traffic")
		}

		lbName := lbRS.Primary.Attributes["load_balancer_name"]
		dnsName := lbRS.Primary.Attributes["dns_name"]
		if lbName == "" || dnsName == "" {
			return fmt.Errorf("load balancer name or dns_name is empty")
		}

		if err := waitForHealthyBackends(t, lbName, lbTrafficBackendCount); err != nil {
			return err
		}

		return checkLBTraffic(t, s, dnsName, lbTrafficHTTPPort, lbTrafficBackendCount)
	}
}

func waitForHealthyBackends(t *testing.T, lbName string, want int) error {
	t.Helper()
	client := testacc.ConfiguredClient.OSC
	if client == nil {
		return fmt.Errorf("configured API client is nil")
	}

	ctx := t.Context()
	start := time.Now()
	attempt := 0
	var lastLog time.Time

	return retry.RetryContext(ctx, lbTrafficHealthTimeout, func() *retry.RetryError {
		attempt++
		resp, err := client.ReadVmsHealth(ctx, osc.ReadVmsHealthRequest{LoadBalancerName: lbName},
			options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			return retry.RetryableError(err)
		}

		entries := resp.BackendVmHealth
		if entries == nil {
			return retry.RetryableError(fmt.Errorf("no backend health entries yet"))
		}

		healthy := 0
		for _, b := range *entries {
			state := ""
			if b.State != nil {
				state = string(*b.State)
			}
			if backendHealthy(state) {
				healthy++
			}
		}

		if len(*entries) < want || healthy < want {
			if attempt == 1 || time.Since(lastLog) >= 30*time.Second {
				t.Logf("waiting for %d healthy backends on %s: %d/%d healthy", want, lbName, healthy, want)
				lastLog = time.Now()
			}
			return retry.RetryableError(fmt.Errorf("%d/%d healthy", healthy, want))
		}

		t.Logf("backends healthy on %s after %s", lbName, time.Since(start).Round(time.Second))
		return nil
	})
}

// checkLBTraffic hits the load balancer, then each VM directly when session affinity hides backends.
func checkLBTraffic(t *testing.T, s *terraform.State, lbHost string, lbPort, backendCount int) error {
	t.Helper()
	ctx := t.Context()
	seen := map[string]int{}

	for i := range lbTrafficRequestCount {
		reqURL := lbTrafficHTTPURL(lbHost, lbPort, fmt.Sprintf("n=%d", i))
		body, err := httpGetRetry(ctx, reqURL)
		if err != nil {
			return fmt.Errorf("LB request %d/%d: %w", i+1, lbTrafficRequestCount, err)
		}
		if !strings.HasPrefix(body, "backend-") {
			return fmt.Errorf("LB request %d/%d: unexpected body %q", i+1, lbTrafficRequestCount, body)
		}
		seen[body]++
	}

	if len(seen) >= backendCount {
		t.Logf("traffic via LB reached %d backend(s): %v", len(seen), seen)
		return nil
	}

	t.Logf("LB saw %d backend(s) from this client (%v); checking VMs directly", len(seen), seen)
	for i := range backendCount {
		vmRS := s.RootModule().Resources[fmt.Sprintf("outscale_vm.lb_traffic_backend.%d", i)]
		if vmRS == nil {
			return fmt.Errorf("not found: outscale_vm.lb_traffic_backend.%d", i)
		}
		ip := vmRS.Primary.Attributes["public_ip"]
		if ip == "" {
			return fmt.Errorf("outscale_vm.lb_traffic_backend.%d: public_ip is empty", i)
		}
		want := fmt.Sprintf("backend-%d", i)
		got, err := httpGetRetry(ctx, lbTrafficHTTPURL(ip, lbPort, ""))
		if err != nil {
			return fmt.Errorf("direct check %s: %w", ip, err)
		}
		if got != want {
			return fmt.Errorf("direct check %s: want %q, got %q", ip, want, got)
		}
	}

	t.Logf("all %d backend VMs serve expected HTTP responses", backendCount)
	return nil
}

func lbTrafficHTTPURL(host string, port int, query string) string {
	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/",
	}
	if query != "" {
		u.RawQuery = query
	}
	return u.String()
}

func backendHealthy(state string) bool {
	switch strings.ToUpper(strings.ReplaceAll(state, "_", "")) {
	case "UP", "INSERVICE", "OK":
		return true
	default:
		return false
	}
}

func httpGetRetry(ctx context.Context, rawURL string) (string, error) {
	var body string
	err := retry.RetryContext(ctx, lbTrafficHTTPRetry, func() *retry.RetryError {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		resp, err := lbTrafficHTTPClient.Do(req)
		if err != nil {
			return retry.RetryableError(err)
		}
		defer func() { _ = resp.Body.Close() }()

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return retry.RetryableError(err)
		}
		body = strings.TrimSpace(string(raw))

		if resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusBadGateway ||
			resp.StatusCode == http.StatusGatewayTimeout {
			return retry.RetryableError(fmt.Errorf("HTTP %d", resp.StatusCode))
		}
		if resp.StatusCode != http.StatusOK {
			return retry.NonRetryableError(fmt.Errorf("HTTP %d: %q", resp.StatusCode, body))
		}
		if body == "" {
			return retry.RetryableError(fmt.Errorf("empty body"))
		}
		return nil
	})
	return body, err
}

func testAccCheckLoadBalancerTrafficDestroy(s *terraform.State) error {
	client := testacc.ConfiguredClient.OSC
	if client == nil {
		return fmt.Errorf("configured API client is nil")
	}

	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		lbName := rs.Primary.ID
		err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			req := osc.ReadLoadBalancersRequest{
				Filters: &osc.FiltersLoadBalancer{
					LoadBalancerNames: &[]string{lbName},
				},
			}

			resp, err := client.ReadLoadBalancers(ctx, req, options.WithRetryTimeout(DefaultTimeout))
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "loadbalancernotfound") {
					return nil
				}
				return retry.RetryableError(err)
			}

			if resp.LoadBalancers == nil {
				return nil
			}
			for _, lb := range *resp.LoadBalancers {
				if lb.LoadBalancerName == lbName {
					return retry.RetryableError(fmt.Errorf("load balancer still exists: %s", lbName))
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccLoadBalancerTrafficConfig(lbName, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_security_group" "sg_lb_traffic" {
  security_group_name = "%[2]s"
  description         = "Used in the terraform acceptance tests"
}

resource "outscale_security_group_rule" "lb_traffic_http" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.sg_lb_traffic.security_group_id
  from_port_range   = %[3]d
  to_port_range     = %[3]d
  ip_protocol       = "tcp"
  ip_range          = "0.0.0.0/0"
}

resource "outscale_vm" "lb_traffic_backend" {
  count              = %[4]d
  image_id           = var.image_id
  vm_type            = var.vm_type
  security_group_ids = [outscale_security_group.sg_lb_traffic.security_group_id]
  user_data = base64encode(<<-EOT
#!/bin/bash
set -euo pipefail
ID="backend-${count.index}"
mkdir -p /var/www/html
echo "$ID" > /var/www/html/index.html
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y lighttpd curl
systemctl enable --now lighttpd
for i in $(seq 1 30); do
  curl -sf "http://127.0.0.1:%[3]d/" | grep -qx "$ID" && exit 0
  sleep 2
done
exit 1
EOT
  )

  lifecycle { ignore_changes = [state] }
}

resource "outscale_load_balancer" "lbu_traffic" {
  load_balancer_name = "%[1]s"
  subregion_names    = [var.subregion]
  secured_cookies    = false
  listeners {
    backend_port           = %[3]d
    backend_protocol       = "HTTP"
    load_balancer_port     = %[3]d
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_load_balancer_attributes" "lbu_traffic" {
  load_balancer_name = outscale_load_balancer.lbu_traffic.load_balancer_name
  health_check {
    healthy_threshold   = 2
    check_interval      = 10
    path                = "/"
    port                = %[3]d
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
`, lbName, sgName, lbTrafficHTTPPort, lbTrafficBackendCount)
}
