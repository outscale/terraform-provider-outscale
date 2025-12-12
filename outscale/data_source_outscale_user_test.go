package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_data_user_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_user.basicTestUser"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
				),
			},
		},
	})
}

const testAccDataUserBasicConfig = `
	resource "outscale_user" "basic_dataUser" {
		user_name = "ACC_user_data1"
		path = "/"
	}
        data "outscale_user" "basicTestUser" {
		filter {
			name = "user_ids"
			values = [outscale_user.basic_dataUser.user_id]
		}
        }
`
