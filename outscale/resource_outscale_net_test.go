package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPINet_Update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPINetConfigUpdateTags("Terraform_net", false),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccOutscaleOAPINetConfigUpdateTags("Terraform_net2", true),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckOutscaleOAPINetExists(n string, res *oscgo.Net) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp oscgo.ReadNetsResponse
		conn := testAccProvider.Meta().(*OutscaleClient)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.OSCAPI.NetApi.ReadNets(context.Background()).ReadNetsRequest(oscgo.ReadNetsRequest{
				Filters: &oscgo.FiltersNet{NetIds: &[]string{rs.Primary.ID}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return err
		}

		if len(resp.GetNets()) != 1 ||
			resp.GetNets()[0].GetNetId() != rs.Primary.ID {
			return fmt.Errorf("Net not found")
		}

		*res = resp.GetNets()[0]

		return nil
	}
}

func testAccOutscaleOAPINetConfigUpdateTags(value string, dhcp bool) string {
	dhcpVal := ""
	if dhcp {
		dhcpVal = "dhcp_options_set_id = outscale_dhcp_option.foo.id"
	}

	return fmt.Sprintf(`
	resource "outscale_dhcp_option" "foo" {
		domain_name         = "test.fr"
		domain_name_servers = ["192.168.12.1"]
	}

	resource "outscale_net" "outscale_net" { 
		ip_range = "10.0.0.0/16"
		%s
		tags { 
			key = "name" 
			value = "%s"
		}
	}
`, dhcpVal, value)
}
