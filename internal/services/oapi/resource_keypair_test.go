package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_keypair_basic(t *testing.T) {
	resourceName := "outscale_keypair.basic_keypair"
	keypairName := acctest.RandomWithPrefix("testacc-keypair")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAcckeypairBasicConfig(keypairName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", keypairName),
				),
			},
		},
	})
}

func TestAccOthers_keypair_Basic_Migration(t *testing.T) {
	keypairName := acctest.RandomWithPrefix("testacc-keypair")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.0.1", testAcckeypairBasicConfig(keypairName)),
	})
}

func TestAccOthers_keypairUpdateTags(t *testing.T) {
	resourceName := "outscale_keypair.update_keypair"
	tagValue1 := "testACC-01"
	tagValue2 := "testACC-02"
	keypairName := acctest.RandomWithPrefix("basic-keypair")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAcckeypairUpdateTags(keypairName, tagValue1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", keypairName),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue1),
				),
			},
			{
				Config: testAcckeypairUpdateTags(keypairName, tagValue2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypair_type"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),

					resource.TestCheckResourceAttr(resourceName, "keypair_name", keypairName),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue2),
				),
			},
		},
	})
}

func TestAccOthers_keypair_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_keypair.update_keypair"
	keypairName := acctest.RandomWithPrefix("basic-keypair")
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-keypair"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			testAcckeypairUpdateTagsWithKey(keypairName, invalidTagKey, tagValue),
			testAcckeypairUpdateTags(keypairName, tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "keypair_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func testAcckeypairBasicConfig(keypair string) string {
	return fmt.Sprintf(`
			resource "outscale_keypair" "basic_keypair" {
				keypair_name = "%s"
			}
		`, keypair)
}

func testAcckeypairUpdateTags(keypairName, value string) string {
	return testAcckeypairUpdateTagsWithKey(keypairName, "name", value)
}

func testAcckeypairUpdateTagsWithKey(keypairName, key, value string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "update_keypair" {
			keypair_name = "%[1]s"
			tags {
				key   = "%[2]s"
				value = "%[3]s"
			}
		}
		`, keypairName, key, value)
}
