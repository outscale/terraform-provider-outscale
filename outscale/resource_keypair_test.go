package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_keypair_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_keypair.basic_keypair"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccKeypairBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", "tesTACC-keypair"),
				),
			},
		},
	})
}

func TestAccOthers_keypairUpdateTags(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_keypair.update_keypair"
	tagValue1 := "testACC-01"
	tagValue2 := "testACC-02"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAcckeypairUpdateTags(tagValue1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", "testTags-keypair"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue1),
				),
			},
			{
				Config: testAcckeypairUpdateTags(tagValue2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", "testTags-keypair"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue2),
				),
			},
		},
	})
}

const testAccKeypairBasicConfig = `
	resource "outscale_keypair" "basic_keypair" {
		keypair_name = "tesTACC-keypair"
	}`

func testAcckeypairUpdateTags(value string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "update_keypair" {
			keypair_name = "testTags-keypair"
			tags {
				key   = "name"
				value = "%s"
		}
		}
	`, value)
}
