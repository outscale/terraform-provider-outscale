package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshotExportTask_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISnapshotExportTaskConfig,
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

var testAccOutscaleOAPISnapshotExportTaskConfig = `
resource "outscale_volume" "test" {
	sub_region = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
}

resource "outscale_snapshot_export_task" "outscale_snapshot_export_task" {
    count = 1

		osu_export {
			disk_image_format = "raw"
			osu_bucket = "test"
		}
    snapshot_id = "${outscale_snapshot.test.id}"
}
`
