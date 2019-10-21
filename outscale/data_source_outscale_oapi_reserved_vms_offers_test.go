package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIReservedVMSOffers(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIReservedVMSOffersConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_reserved_vms_offers.test", "reserved_instances_offerings_set"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIReservedVMSOffersConfig = `
data "outscale_reserved_vms_offers" "test" {}
`
