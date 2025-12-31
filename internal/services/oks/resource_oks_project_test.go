package oks_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOKSProject_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("project-basic")
	resourceName := "outscale_oks_project.project"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: oksProjectConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "10.50.0.0/18"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
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
