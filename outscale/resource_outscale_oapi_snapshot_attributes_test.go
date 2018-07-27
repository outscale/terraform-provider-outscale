package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshotAttributes_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var snapshotID string
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISnapshotAttributesConfig(true, accountID),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_snapshot.test", "id", &snapshotID),
					testAccOutscaleOAPISnapshotAttributesExists(&accountID, &snapshotID),
				),
			},
			// Drop just create volume permission to test destruction
			resource.TestStep{
				Config: testAccOutscaleOAPISnapshotAttributesConfig(false, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleOAPISnapshotAttributesDestroyed(&accountID, &snapshotID),
				),
			},
		},
	})
}

func testAccOutscaleOAPISnapshotAttributesExists(accountID, snapshotID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		if has, err := hasOAPICreateVolumePermission(conn, *snapshotID, *accountID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("create volume permission does not exist for '%s' on '%s'", *accountID, *snapshotID)
		}
		return nil
	}
}

func testAccOutscaleOAPISnapshotAttributesDestroyed(accountID, snapshotID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		if has, err := hasOAPICreateVolumePermission(conn, *snapshotID, *accountID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("create volume permission still exists for '%s' on '%s'", *accountID, *snapshotID)
		}
		return nil
	}
}

func testAccOutscaleOAPISnapshotAttributesConfig(includeCreateVolumePermission bool, aid string) string {
	base := `
resource "outscale_volume" "description_test" {
	sub_region = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.description_test.id}"
	description = "Snapshot Acceptance Test"
}
`

	if !includeCreateVolumePermission {
		return base
	}

	return base + fmt.Sprintf(`
resource "outscale_snapshot_attributes" "self-test" {
	snapshot_id = "${outscale_snapshot.test.id}"
	permission_to_create_volume = [{
		create = [{
			account_id = "%s"
		}]
	}]
}
`, aid)
}
