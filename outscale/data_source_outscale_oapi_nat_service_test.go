package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPINatServiceDataSource_Instance(t *testing.T) {
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
				Config: testAccCheckOutscaleOAPINatServiceDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleNatServiceDataSourceID("data.outscale_nat_service.nat"),
					resource.TestCheckResourceAttr("data.outscale_nat_service.nat", "subnet_id", "subnet-861fbecc"),
				),
			},
		},
	})
}

const testAccCheckOutscaleOAPINatServiceDataSourceConfig = `
data "outscale_nat_service" "nat" {
	nat_service_id = "nat-08f41400"
}
`
