package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOutscaleOAPILin_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.Net

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckOutscaleLinDestroyed, // we need to create the destroyed test case
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinExists("outscale_net.vpc", &conf),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func TestAccNet_UpdateTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILinConfigUpdateTags("Terraform_net"),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccOutscaleOAPILinConfigUpdateTags("Terraform_net2"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckOutscaleOAPILinExists(n string, res *oscgo.Net) resource.TestCheckFunc {
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

const testAccOutscaleOAPILinConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name" 
			value = "testacc-net-rs"
		}	
	}
`

func testAccOutscaleOAPILinConfigUpdateTags(value string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" { 
		ip_range = "10.0.0.0/16"
		tags { 
			key = "name" 
			value = "%s"
		}
	   }
`, value)
}
