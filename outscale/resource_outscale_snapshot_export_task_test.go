package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_SnapshotExportTask_basic(t *testing.T) {
	osuBucketNames := []string{
		acctest.RandomWithPrefix("terraform-export-bucket-"),
		acctest.RandomWithPrefix("terraform-export-bucket-"),
	}
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
				Config: testAccOutscaleOAPISnapshotExportTaskConfig("", osuBucketNames[0], utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotExportTaskExists("outscale_snapshot_export_task.outscale_snapshot_export_task"),
				),
			},
			{
				Config: testAccOutscaleOAPISnapshotExportTaskConfig(tags, osuBucketNames[1], utils.GetRegion()),
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

func testAccOutscaleOAPISnapshotExportTaskConfig(tags, osuBucketName, region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume_snap" {
    subregion_name   = "%[3]sa"
    size                = 10
}
resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume_snap.volume_id
}
resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
	snapshot_id                     = outscale_snapshot.outscale_snapshot.snapshot_id
	osu_export {
		disk_image_format = "qcow2"
        osu_bucket        = "%[2]s"
        osu_prefix        = "new-export"
	}
	%[1]s
}
	`, tags, osuBucketName, region)
}
