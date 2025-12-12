package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_Bacic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configNetBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_net.basic_net", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func TestAccNet_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.0.1", configNetBasic),
	})
}

func TestAccNet_UpdateTags(t *testing.T) {
	netName := "outscale_net.basic_net"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configNetUpdateTags("NetTags"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(netName, "tags.#"),
				),
			},
			{
				Config: configNetUpdateTags("NetTags2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(netName, "tags.#"),
				),
			},
		},
	})
}

const configNetBasic = `
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-net-rs"
		}
	}
`

func configNetUpdateTags(tagValue string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "name"
			value = "%s"
		}
	   }
`, tagValue)
}
