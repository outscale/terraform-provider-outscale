package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNet_Attr_import(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_net_attributes.outscale_net_attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPILinAttrConfig,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}
