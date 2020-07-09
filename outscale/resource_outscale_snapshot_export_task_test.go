package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPISnapshotExportTask_basic(t *testing.T) {
	tags := `tags {
		key = "test"
		value = "test"
	}
	tags {
		key = "test-1"
		value = "test-1"
	}`
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotExportTaskConfig(""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotExportTaskExists("outscale_snapshot_export_task.outscale_snapshot_export_task"),
				),
			},
			{
				Config: testAccOutscaleOAPISnapshotExportTaskConfig(tags),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotExportTaskExists("outscale_snapshot_export_task.outscale_snapshot_export_task"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISnapshotExportTaskExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No image task id is set")
		}

		return nil
	}
}

func testAccOutscaleOAPISnapshotExportTaskConfig(tags string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume_snap" {
    subregion_name   = "eu-west-2a"
    size                = 10
}
resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}
resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
	snapshot_id                     = outscale_snapshot.outscale_snapshot.snapshot_id
	osu_export {
		disk_image_format = "qcow2"
        osu_bucket        = "terraform-export-bucket"
        osu_prefix        = "new-export"
	}
	%s
}
	`, tags)
}
