package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPILinAttr_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinAttrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_lin_attributes.outscale_lin_attributes", "dns_support_enabled", "true"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinAttrConfig = `
resource "outscale_lin_attributes" "outscale_lin_attributes" {
	dns_support_enabled = true
	lin_id = "vpc-5b79bc69"
	attribute            = "enableDnsSupport"
}
`
