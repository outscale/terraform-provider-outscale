package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPINatServicesDataSource_Instance(t *testing.T) {
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
			{
				Config: testAccCheckOutscaleOAPINatServicesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNatServiceDataSourceID("data.outscale_nat_services.nat"),
					resource.TestCheckResourceAttr("data.outscale_nat_services.nat", "nat_service.#", "1"),
					resource.TestCheckResourceAttr("data.outscale_nat_services.nat", "nat_service.0.subnet_id", "subnet-861fbecc"),
				),
			},
		},
	})
}

const testAccCheckOutscaleOAPINatServicesDataSourceConfig = `
data "outscale_nat_services" "nat" {
	filter {
		name = "nat_service_ids" 
		values = ["nat-08f41400"]
	}
}
`
