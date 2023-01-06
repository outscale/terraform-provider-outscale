package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOutscaleOAPISnapshotAttributes_Basic(t *testing.T) {
	t.Skip()
	var snapshotID string
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotAttributesConfig(true, false, accountID),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_snapshot.test", "id", &snapshotID),
				),
			},
			{
				Config: testAccOutscaleOAPISnapshotAttributesConfig(true, true, accountID),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_snapshot.test", "id", &snapshotID),
				),
			},
		},
	})
}

func testAccOutscaleOAPISnapshotAttributesConfig(includeAddition, includeRemoval bool, aid string) string {
	base := fmt.Sprintf(`
		resource "outscale_volume" "description_test" {
			subregion_name = "eu-west-2a"
			size           = 1
		}
		
		resource "outscale_snapshot" "test" {
			volume_id   = "${outscale_volume.description_test.id}"
			description = "Snapshot Acceptance Test"
		}
		
		resource "outscale_snapshot_attributes" "self-test" {
			snapshot_id = "${outscale_snapshot.test.id}"
		
			permissions_to_create_volume_removals {
				account_ids = ["%s"]
			}
		}
	`, aid)

	if includeAddition {
		return base + fmt.Sprintf(`
			resource "outscale_snapshot_attributes" "additions" {
				snapshot_id = "${outscale_snapshot.test.id}"
			
				permissions_to_create_volume_additions {
					account_ids = ["%s"]
				}
			}
		`, aid)
	}

	if includeRemoval {
		return base + fmt.Sprintf(`
		resource "outscale_snapshot_attributes" "removals" {
			snapshot_id = "${outscale_snapshot.test.id}"
		
			permissions_to_create_volume_removals {
				account_ids = ["%s"]
			}
		}
		`, aid)
	}
	return base
}
