package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscalePublicIPLink_importBasic(t *testing.T) {
	resourceName := "outscale_public_ip_link.bar"

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscalePublicIPAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePublicIPAssociationConfig,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"associate_public_ip_address", "user_data", "security_group"},
			},
		},
	})
}
