package outscale

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPISnapshotExportTaskDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotExportTaskDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotExportTaskDataSourceID("data.outscale_snapshot_export_task.export_task"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSnapshotExportTaskDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Snapshot Export Task data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot Export Task data source ID not set")
		}
		return nil
	}
}

var testAccOutscaleOAPISnapshotExportTaskDataSourceConfig = `
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
}

data "outscale_snapshot_export_task" "export_task" {
	filter {
		name = "task_ids"
		values = [outscale_snapshot_export_task.outscale_snapshot_export_task.id]
	}
}
`
