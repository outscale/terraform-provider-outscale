package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPINatService_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	var natGateway oapi.NatService

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nat_service.gateway",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPINatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPINatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPINatGatewayExists("outscale_nat_service.nat_service", &natGateway),
				),
			},
		},
	})
}

func testAccCheckOAPINatGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nat_service" {
			continue
		}

		// Try to find the resource

		var resp *oapi.POST_ReadNatServicesResponses

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadNatServices(oapi.ReadNatServicesRequest{
				Filters: oapi.FiltersNatService{NatServiceIds: []string{rs.Primary.ID}},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				if strings.Contains(err.Error(), "NatGatewayNotFound:") {
					return nil
				}
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("[DEBUG] Error reading Nat Service (%s)", errString)
		}

		response := resp.OK

		if err == nil {
			status := map[string]bool{
				"deleted":  true,
				"deleting": true,
				"failed":   true,
			}
			if _, ok := status[strings.ToLower(response.NatServices[0].State)]; len(response.NatServices) > 0 && !ok {
				return fmt.Errorf("still exists")
			}

			return nil
		}

	}

	return nil
}

func testAccCheckOAPINatGatewayExists(n string, ng *oapi.NatService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		var resp *oapi.POST_ReadNatServicesResponses

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadNatServices(oapi.ReadNatServicesRequest{
				Filters: oapi.FiltersNatService{NatServiceIds: []string{rs.Primary.ID}},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("[DEBUG] Error reading Nat Service (%s)", errString)
		}

		response := resp.OK

		if len(response.NatServices) == 0 {
			return fmt.Errorf("NatGateway not found")
		}

		*ng = response.NatServices[0]

		return nil
	}
}

const testAccOAPINatGatewayConfig = `
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    net_id = "${outscale_net.outscale_net.net_id}"
    ip_range = "10.0.0.0/18"
}

resource "outscale_public_ip" "outscale_public_ip" {
}

resource "outscale_nat_service" "nat_service" {
   subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
   public_ip_id = "${outscale_public_ip.outscale_public_ip.public_ip_id}"
}
`
