package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_users_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_users.basicTestUsers"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleDataUserBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "users.#"),
				),
			},
		},
	})
}

const testAccOutscaleDataUserBasicConfig = `
	resource "outscale_user" "basic_data_users" {
	  user_name = "ACC_test_data1"
	  path = "/"
	}
        data "outscale_users" "basicTestUsers" {
          depends_on = [outscale_user.basic_data_users]
        }`
