package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"
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
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
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
