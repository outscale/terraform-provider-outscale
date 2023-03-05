package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPISubNet_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.Subnet

	resourceName := "outscale_subnet.subnet"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISubNetDestroyed, // we need to create the destroyed test case
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISubnetConfig(utils.GetRegion(), false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubNetExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "map_public_ip_on_launch", "false"),
				),
			},
			{
				Config: testAccOutscaleOAPISubnetConfig(utils.GetRegion(), true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubNetExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "map_public_ip_on_launch", "true"),
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
			rp, httpResp, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(oscgo.ReadSubnetsRequest{
				Filters: &oscgo.FiltersSubnet{
					SubnetIds: &[]string{rs.Primary.ID},
				},
			}).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
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
			rp, httpResp, err := conn.SubnetApi.ReadSubnets(context.Background()).ReadSubnetsRequest(oscgo.ReadSubnetsRequest{
				Filters: &oscgo.FiltersSubnet{
					SubnetIds: &[]string{rs.Primary.ID},
				},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
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

func testAccOutscaleOAPISubnetConfig(region string, mapPublicIpOnLaunch bool) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-subnet-rs"
			}
		}

		resource "outscale_subnet" "subnet" {
			ip_range       = "10.0.0.0/24"
			subregion_name = "%sb"
			net_id         = outscale_net.net.id
			map_public_ip_on_launch = %v
			tags {
				key   = "name"
				value = "terraform-subnet"
			}
		}
	`, region, mapPublicIpOnLaunch)
}
