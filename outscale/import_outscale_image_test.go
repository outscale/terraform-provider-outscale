package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleImage_importBasic(t *testing.T) {
	resourceName := "outscale_image.foo"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccImageConfig_basic(rInt),
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
