package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleLinInternetGatewaysDatasource_basic(t *testing.T) {
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
				Config: testAccOutscaleLinInternetGatewaysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_lin_internet_gateways.outscale_lin_internet_gateways", "internet_gateway_set.#", "1"),
				),
			},
		},
	})
}

const testAccOutscaleLinInternetGatewaysDatasourceConfig = `
resource "outscale_lin_internet_gateway" "gateway" {}

data "outscale_lin_internet_gateways" "outscale_lin_internet_gateways" {
  filter {
		name = "internet-gateway-id"
		values = ["${outscale_lin_internet_gateway.gateway.id}"]
	}
}
`
