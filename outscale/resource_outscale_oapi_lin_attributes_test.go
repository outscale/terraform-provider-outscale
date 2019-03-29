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

	if !oapi {
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
						"outscale_net_attributes.outscale_net_attributes", "dns_support_enabled", "true"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinAttrConfig = `

resource "outscale_net" "vpc" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_net_attributes" "outscale_net_attributes" {
	net_id = "${outscale_net.vpc.id}"
	dhcp_options_set_id = "set-id"
}
`
