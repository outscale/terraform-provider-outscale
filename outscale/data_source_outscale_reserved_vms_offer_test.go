package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIReservedVMSOffer(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIReservedVMSOfferConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_reserved_vms_offer.test", "reserved_instances_offering_id"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIReservedVMSOfferConfig = `
	data "outscale_reserved_vms_offer" "test" {}
`
