package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOKSProject_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("project-basic")
	resourceName := "outscale_oks_project.project"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: oksProjectConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "10.50.0.0/18"),
				),
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

func oksProjectConfig(name string) string {
	return fmt.Sprintf(`
		resource "outscale_oks_project" "project" {
			name = "%s"
			cidr = "10.50.0.0/18"
			region = "eu-west-2"
			disable_api_termination = false

			tags = {
				test = "TestAccProjectBasic"
			}
		}
`, name)
}
