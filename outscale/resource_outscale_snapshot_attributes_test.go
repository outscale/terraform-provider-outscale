package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleSnapshotAttributes_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var snapshotID string
	accountID := os.Getenv("OUTSCALE_ACCOUNT")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleSnapshotAttributesConfig(true, accountID),
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceGetAttr("outscale_snapshot.test", "id", &snapshotID),
					testAccOutscaleSnapshotAttributesExists(&accountID, &snapshotID),
				),
			},
			// Drop just create volume permission to test destruction
			resource.TestStep{
				Config: testAccOutscaleSnapshotAttributesConfig(false, accountID),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleSnapshotAttributesDestroyed(&accountID, &snapshotID),
				),
			},
		},
	})
}

func testAccOutscaleSnapshotAttributesExists(accountID, snapshotID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		if has, err := hasCreateVolumePermission(conn, *snapshotID, *accountID); err != nil {
			return err
		} else if !has {
			return fmt.Errorf("create volume permission does not exist for '%s' on '%s'", *accountID, *snapshotID)
		}
		return nil
	}
}

func testAccOutscaleSnapshotAttributesDestroyed(accountID, snapshotID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		if has, err := hasCreateVolumePermission(conn, *snapshotID, *accountID); err != nil {
			return err
		} else if has {
			return fmt.Errorf("create volume permission still exists for '%s' on '%s'", *accountID, *snapshotID)
		}
		return nil
	}
}

func testAccOutscaleSnapshotAttributesConfig(includeCreateVolumePermission bool, aid string) string {
	base := `
resource "outscale_volume" "description_test" {
	availability_zone = "eu-west-2a"
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
	create_volume_permission = [{
		add = [{
			user_id = "%s"
		}]
	}]
}
`, aid)
}
