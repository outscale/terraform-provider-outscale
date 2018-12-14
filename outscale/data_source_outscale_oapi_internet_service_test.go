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
					testAccCheckState("data.outscale_internet_service.outscale_internet_serviced"),
					resource.TestCheckResourceAttrSet("data.outscale_internet_service.outscale_internet_serviced", "internet_service_id"),
				),
			},
		},
	})
}

const testAccOutscaleOAPIINternetServiceDatasourceConfig = `
resource "outscale_internet_service" "outscale_internet_service" {}

data "outscale_internet_service" "outscale_internet_serviced" {
	filter {
		name = "internet_service_ids"
		values = ["${outscale_internet_service.outscale_internet_service.internet_service_id}"]
	}
}
`
