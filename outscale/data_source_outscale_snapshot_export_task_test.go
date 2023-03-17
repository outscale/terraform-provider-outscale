package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SnapshotExportTask_DataSource(t *testing.T) {
	t.Parallel()
	imageName := acctest.RandomWithPrefix("terraform-export-")
	region := os.Getenv("OUTSCALE_REGION")
	dataSourceName := "data.outscale_snapshot_export_task.export_task"
	dataSourcesName := "data.outscale_snapshot_export_tasks.export_tasks"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_SnapshotExportTask_DataSource_Config(region, imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "snapshot_export_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "task_id"),
				),
			},
		},
	})
}

func testAcc_SnapshotExportTask_DataSource_Config(region, testName string) string {
	var stringTemplate = `
		resource "outscale_volume" "outscale_volume_snap" {
			subregion_name   = "%[1]a"
			size             = 1
		}

		resource "outscale_snapshot" "outscale_snapshot" {
			volume_id = outscale_volume.outscale_volume_snap.volume_id
		}

		resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
			snapshot_id           = outscale_snapshot.outscale_snapshot.snapshot_id
			osu_export {
				disk_image_format = "qcow2"
				osu_bucket        = "%[2]s"
				osu_prefix        = "new-export"
			}
		}

		data "outscale_snapshot_export_task" "export_task" {
			filter {
				name = "task_ids"
				values = [outscale_snapshot_export_task.outscale_snapshot_export_task.id]
			}
		}

		data "outscale_snapshot_export_tasks" "export_tasks" {
			filter {
				name = "task_ids"
				values = [outscale_snapshot_export_task.outscale_snapshot_export_task.id]
			}
		}
		`
	return fmt.Sprintf(stringTemplate, region, testName)
}
