package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPILinInternetGatewayDatasource_basic(t *testing.T) {
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
				Config: testAccOutscaleOAPILinInternetGatewayDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_lin_internet_gateway", "test.lin_to_lin_internet_gateway_link.#", "1"),
				),
			},
		},
	})
}

const testAccOutscaleOAPILinInternetGatewayDatasourceConfig = `
resource "outscale_lin_internet_gateway" "gateway" {}

data "outscale_lin_internet_gateway" "test" {
	filter {
		name = "internet-gateway-id"
		values = ["${outscale_lin_internet_gateway.gateway.id}"]
	}
}
`
