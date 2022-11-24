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

func TestAccLin_basic(t *testing.T) {
	t.Parallel()
	var conf1 oscgo.Net
	var conf2 oscgo.Net

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckLinDestroyed, // we need to create the destroyed test case
		Steps: []resource.TestStep{
			{
				Config: testAccLinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinExists("outscale_net.vpc.0", &conf1),
					testAccCheckLinExists("outscale_net.vpc.1", &conf2),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc.0", "ip_range", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc.1", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func TestAccLin_UpdateTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLinConfigUpdateTags("Terraform_net"),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccLinConfigUpdateTags("Terraform_net2"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckLinExists(n string, res *oscgo.Net) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp oscgo.ReadNetsResponse
		conn := testAccProvider.Meta().(*Client)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.OSCAPI.NetApi.ReadNets(context.Background()).ReadNetsRequest(oscgo.ReadNetsRequest{
				Filters: &oscgo.FiltersNet{NetIds: &[]string{rs.Primary.ID}},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
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

const testAccLinConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"
		count = 2

		tags {
			key = "Name" 
			value = "testacc-net-rs"
		}	
	}
`

func testAccLinConfigUpdateTags(value string) string {
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
