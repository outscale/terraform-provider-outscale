package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIINternetServiceDatasource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIINternetServiceDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_internet_service", "test.net_to_internet_service_link.#", "1"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIINternetServiceDatasourceConfig = `
resource "outscale_internet_service" "gateway" {}

data "outscale_internet_service" "test" {
	filter {
		name = "internet-gateway-id"
		values = ["${outscale_internet_service.gateway.id}"]
	}
}
`
