package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDSOutscaleUserAPIKey_basic(t *testing.T) {
	rName := fmt.Sprintf("test-user-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDSOutscaleUserAPIKeyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("outscale_user_api_keys.a_key", "secret_access_key"),
					resource.TestCheckResourceAttr("outscale_user_api_keys.a_key", "user_name", rName),
					resource.TestCheckResourceAttr("data.outscale_user_api_keys.test_key", "user_name", rName),
					resource.TestCheckResourceAttr("data.outscale_user_api_keys.test_key", "access_key_metadata.#", "1"),
					resource.TestCheckResourceAttr("data.outscale_user_api_keys.test_key", "access_key_metadata.0.user_name", rName),
				),
			},
		},
	})
}

func testAccDSOutscaleUserAPIKeyConfig(rName string) string {
	return fmt.Sprintf(`
resource "outscale_user" "a_user" {
        user_name = "%s"
}
resource "outscale_user_api_keys" "a_key" {
        user_name = "${outscale_user.a_user.user_name}"
}

data "outscale_user_api_keys" "test_key" {
        user_name = "${outscale_user_api_keys.a_key.user_name}"
}
`, rName)
}
