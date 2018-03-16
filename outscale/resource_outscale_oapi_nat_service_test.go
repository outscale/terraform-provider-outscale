package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPINatService_basic(t *testing.T) {
	var natGateway fcu.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nat_service.gateway",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOAPINatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPINatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPINatGatewayExists("outscale_nat_service.gateway", &natGateway),
				),
			},
		},
	})
}

func testAccCheckOAPINatGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nat_service" {
			continue
		}

		// Try to find the resource

		var resp *fcu.DescribeNatGatewaysOutput

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.DescribeNatGateways(&fcu.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err == nil {
			status := map[string]bool{
				"deleted":  true,
				"deleting": true,
				"failed":   true,
			}
			if _, ok := status[strings.ToLower(*resp.NatGateways[0].State)]; len(resp.NatGateways) > 0 && !ok {
				return fmt.Errorf("still exists")
			}

			return nil
		}
		if err != nil {
			if strings.Contains(err.Error(), "NatGatewayNotFound:") {
				return nil
			}

			return err
		}

	}

	return nil
}

func testAccCheckOAPINatGatewayExists(n string, ng *fcu.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		var resp *fcu.DescribeNatGatewaysOutput
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error

			resp, err = conn.VM.DescribeNatGateways(&fcu.DescribeNatGatewaysInput{
				NatGatewayIds: []*string{aws.String(rs.Primary.ID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}
		if len(resp.NatGateways) == 0 {
			return fmt.Errorf("NatGateway not found")
		}

		*ng = *resp.NatGateways[0]

		return nil
	}
}

const testAccOAPINatGatewayConfig = `
resource "outscale_nat_service" "gateway" {
    reservation_id = "eipalloc-32e506e8"
    subnet_id = "subnet-861fbecc"
}
`
