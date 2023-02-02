package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_SecurityGroup_DataSource(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	dataSourceName := "data.outscale_security_group.sg"
	dataSourcesName := "data.outscale_security_groups.sgs"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_SecurityGroup_DataSource_Config(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcesName, "security_groups.#", "2"),

					resource.TestCheckResourceAttr(dataSourceName, "description", "Used in the terraform acceptance tests datasource"),
				),
			},
		},
	})
}

func testAcc_SecurityGroup_DataSource_Config(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "test" {
			description = "Used in the terraform acceptance tests datasource"
			security_group_name = "test-1-%[1]d"
			tag = {
				Name = "tf-acctest"
				Seed = "%[1]d"
			}
		}

		resource "outscale_security_group" "test2" {
			description = "Used in the terraform acceptance tests datasource"
			security_group_name = "test-2-%[1]d"
			tag = {
				Name = "tf-acctest"
				Seed = "%[1]d"
			}
		}

		data "outscale_security_group" "sg" {
			filter {
				name = "security_group_ids"
				values = ["${outscale_security_group.test.security_group_id}"]
			}
		}
		
		data "outscale_security_groups" "sgs" {
			filter {
				name = "security_group_ids"
				values = ["${outscale_security_group.test.security_group_id}", "${outscale_security_group.test2.security_group_id}"]
			}
		}`,
		rInt)
}
