package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPIInternetServicesDatasource_basic(t *testing.T) {
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
				Config: testAccOutscaleOAPIInternetServicesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_internet_services.outscale_internet_services", "internet_services.#", "1"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIInternetServicesDatasourceConfig = `
resource "outscale_internet_service" "gateway" {}

data "outscale_internet_services" "outscale_internet_services" {
  filter {
		name = "internet_service_id"
		values = ["${outscale_internet_service.gateway.id}"]
	}
}
`
