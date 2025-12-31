package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_data_user_basic(t *testing.T) {
	resourceName := "data.outscale_user.basicTestUser"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testacc.PreCheck(t) },
		Providers: testacc.SDKProviders,
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
