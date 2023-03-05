package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPISnapshotExportTaskDataSource_basic(t *testing.T) {
	t.Parallel()
	imageName := acctest.RandomWithPrefix("terraform-export-")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotExportTaskDataSourceConfig(imageName, utils.GetRegion()),
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

func testAccOutscaleOAPISnapshotExportTaskDataSourceConfig(testName, region string) string {
	var stringTemplate = `
		resource "outscale_volume" "outscale_volume_snap" {
			subregion_name   = "%[2]sa"
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

		data "outscale_snapshot_export_task" "export_task" {
			filter {
				name = "task_ids"
				values = [outscale_snapshot_export_task.outscale_snapshot_export_task.id]
			}
		}
		`
	return fmt.Sprintf(stringTemplate, testName)
}
