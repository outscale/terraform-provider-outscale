package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshotCopy_Basic(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleSnapshotCopyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccOutscaleSnapshotCopyExists("outscale_snapshot_import"),
				),
			},
		},
	})
}

func testAccOutscaleSnapshotCopyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccOutscaleSnapshotImportConfig() string {
	return fmt.Sprintf(`
resource "outscale_snapshot_import" "test" {
	snapshot_location = ""
snapshot_size = ""
}
`)
}

func testAccOutscaleSnapshotCopyConfig() string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" {
	availability_zone = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
	description = "Snapshot Acceptance Test"
}

resource "outscale_snapshot_copy" "test" {
	source_region =  "eu-west-2"
	source_snapshot_id = "${outscale_snapshot.test.id}"
}
`)
}
