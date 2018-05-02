package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleReservedVMSOffers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleReservedVMSOffersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_reserved_vms_offers.test", "reserved_instances_offerings_set"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleReservedVMSOffersConfig = `
data "outscale_reserved_vms_offers" "test" {}
`
