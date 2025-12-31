package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_users_basic(t *testing.T) {
	resourceName := "data.outscale_users.basicTestUsers"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
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
