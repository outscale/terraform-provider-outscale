package outscale

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func TestAccOutscaleOAPISubNet_basic(t *testing.T) {
	var conf oscgo.Subnet
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISubNetDestroyed, // we need to create the destroyed test case
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISubnetConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubNetExists("outscale_subnet.subnet", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISubNetExists(n string, res *oscgo.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet id is set")
		}
		var resp oscgo.ReadSubnetsResponse
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			var err error
			resp, _, err = conn.SubnetApi.ReadSubnets(context.Background(), &oscgo.ReadSubnetsOpts{
				ReadSubnetsRequest: optional.NewInterface(oscgo.ReadSubnetsRequest{
					Filters: &oscgo.FiltersSubnet{
						SubnetIds: &[]string{rs.Primary.ID},
					},
				}),
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
			return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", err)
		}

		if len(resp.GetSubnets()) != 1 ||
			resp.GetSubnets()[0].GetSubnetId() != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*res = resp.GetSubnets()[0]

		return nil
	}
}

func testAccCheckOutscaleOAPISubNetDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_subnet" {
			continue
		}

		// Try to find a subnet
		var resp oscgo.ReadSubnetsResponse
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		err := resource.Retry(30*time.Second, func() *resource.RetryError {
			var err error
			resp, _, err = conn.SubnetApi.ReadSubnets(context.Background(), &oscgo.ReadSubnetsOpts{
				ReadSubnetsRequest: optional.NewInterface(oscgo.ReadSubnetsRequest{
					Filters: &oscgo.FiltersSubnet{
						SubnetIds: &[]string{rs.Primary.ID},
					},
				}),
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
			if len(resp.GetSubnets()) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		// Verify the error is what we want
		if !strings.Contains(fmt.Sprintf("%s", err), "InvalidSubnet.NotFound") {
			return err
		}
	}
	return nil
}

func testAccOutscaleOAPISubnetConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-subnet-rs"
			}
		}

		resource "outscale_subnet" "subnet" {
			ip_range       = "10.0.0.0/16"
			subregion_name = "%sa"
			net_id         = "${outscale_net.net.id}"

			tags {
				key   = "name"
				value = "terraform-subnet"
			}
		}
	`, region)
}
