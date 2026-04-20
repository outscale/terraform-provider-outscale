package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_SnapshotAttributes_Basic(t *testing.T) {
	accountID := os.Getenv("OUTSCALE_ACCOUNT")
	resourceName := "outscale_snapshot.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotAttributesConfig(true, false, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume.#", "0"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume.#", "1"),
				),
			},
			{
				Config: testAccOutscaleSnapshotAttributesConfig(false, true, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume.#", "1"),
				),
			},
		},
	})
}

func TestAccOthers_SnapshotAttributes_Migration(t *testing.T) {
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			testAccOutscaleSnapshotAttributesConfig(true, false, accountID),
		),
	})
}

func testAccOutscaleSnapshotAttributesConfig(includeAddition, includeRemoval bool, accountID string) string {
	base := fmt.Sprintf(`
		resource "outscale_volume" "vol" {
			subregion_name = "%[2]sa"
			size           = 1
		}

		resource "outscale_snapshot" "test" {
			volume_id   = outscale_volume.vol.id
			description = "Snapshot Acceptance Test"
		}
	`, accountID, utils.GetRegion())

	if includeAddition {
		return base + fmt.Sprintf(`
			resource "outscale_snapshot_attributes" "additions" {
				snapshot_id = outscale_snapshot.test.id

				permissions_to_create_volume_additions {
					account_ids = ["%s"]
				}
			}
		`, accountID)
	}

	if includeRemoval {
		return base + fmt.Sprintf(`
		resource "outscale_snapshot_attributes" "removals" {
			snapshot_id = outscale_snapshot.test.id

			permissions_to_create_volume_removals {
				account_ids = ["%s"]
			}
		}
		`, accountID)
	}
	return base
}
