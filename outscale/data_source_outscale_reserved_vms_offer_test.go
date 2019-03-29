package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleReservedVMSOffer(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	t.Skip()

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
				Config: testAccDataSourceOutscaleReservedVMSOfferConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_reserved_vms_offer.test", "reserved_instances_offering_id"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleReservedVMSOfferConfig = `
data "outscale_reserved_vms_offer" "test" {}
`
