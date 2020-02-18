package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleLinAttr_import(t *testing.T) {
	resourceName := "outscale_lin_attributes.outscale_lin_attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinAttrConfig,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}
