package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_data_user_group_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_user_group.basicUTest"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_group_name"),
				),
			},
		},
	})
}

const testAccDataUserGroupBasicConfig = `
	resource "outscale_user_group" "uGData" {
		user_group_name = "TestACC_udata"
		path = "/"
	}
	data "outscale_user_group" "basicUTest" {
	    user_group_name = outscale_user_group.uGData.user_group_name
    }
`
