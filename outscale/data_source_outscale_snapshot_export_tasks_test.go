package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPISnapshotExportTasksDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotExportTasksDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotExportTaskDataSourceID("data.outscale_snapshot_export_tasks.export_tasks"),
				),
			},
		},
	})
}

var testAccOutscaleOAPISnapshotExportTasksDataSourceConfig = `
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

data "outscale_snapshot_export_tasks" "export_tasks" {
	filter {
		name = "task_ids"
		values = [outscale_snapshot_export_task.outscale_snapshot_export_task.id]
	}
}
`
