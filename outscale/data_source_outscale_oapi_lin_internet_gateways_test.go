package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPILinInternetGatewaysDatasource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinInternetGatewaysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_lin_internet_gateways.outscale_lin_internet_gateways", "lin_internet_gateway.#", "1"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinInternetGatewaysDatasourceConfig = `
resource "outscale_lin_internet_gateway" "gateway" {}

data "outscale_lin_internet_gateways" "outscale_lin_internet_gateways" {
  filter {
		name = "lin_internet_gateway_id"
		values = ["${outscale_lin_internet_gateway.gateway.id}"]
	}
}
`
