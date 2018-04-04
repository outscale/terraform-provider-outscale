package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleLinInternetGatewayDatasource_basic(t *testing.T) {
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
				Config: testAccOutscaleLinInternetGatewayDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLinIGDataSourceID("data.outscale_lin_internet_gateway.test"),
					testAccCheckOutscaleLinIGDataSourceID("data.outscale_lin_internet_gateway.by_id"),
				),
			},
		},
	})
}

func testAccCheckOutscaleLinIGDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Lin Internet Gateway data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Lin Internet Gateway data source ID not set")
		}
		return nil
	}
}

const testAccOutscaleLinInternetGatewayDatasourceConfig = `
resource "outscale_lin_internet_gateway" "gateway" {}

data "outscale_lin_internet_gateway" "test" {
	filter {
		name = "internet-gateway-id"
		values = ["${outscale_lin_internet_gateway.gateway.id}"]
	}
}
data "outscale_lin_internet_gateway" "by_id" {
	internet_gateway_id = "${outscale_lin_internet_gateway.gateway.id}"
}
`
