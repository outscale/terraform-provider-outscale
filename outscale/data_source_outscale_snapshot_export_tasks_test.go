package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOutscaleOAPISnapshotExportTasksDataSource_basic(t *testing.T) {
	imageName := acctest.RandomWithPrefix("terraform-export")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotExportTasksDataSourceConfig(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotExportTaskDataSourceID("data.outscale_snapshot_export_tasks.export_tasks"),
				),
			},
		},
	})
}

func testAccOutscaleOAPISnapshotExportTasksDataSourceConfig(testName string) string {
	var stringTemplate = `
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
					osu_bucket        = "%s"
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
	return fmt.Sprintf(stringTemplate, testName)
}
