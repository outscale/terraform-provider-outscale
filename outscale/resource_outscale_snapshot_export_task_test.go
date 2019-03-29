package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleSnapshotExportTask_basic(t *testing.T) {
	rInt := acctest.RandIntRange(64512, 65534)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleSnapshotExportTaskConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotExportTaskExists("outscale_snapshot_export_tasks.outscale_snapshot_export_tasks"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSnapshotExportTaskExists(n string) resource.TestCheckFunc {
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

func testAccOutscaleSnapshotExportTaskConfig(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" { 
  availability_zone = "eu-west-2a" 
  size = 1 
}

resource "outscale_snapshot" "test" { 
  volume_id = "${outscale_volume.test.id}" 
}

resource "outscale_snapshot_export_tasks" "outscale_snapshot_export_tasks" {
  snapshot_id = "${outscale_snapshot.test.id}"

	export_to_osu_disk_image_format = "raw"
  export_to_osu_bucket = "customer_tooling_%d"
 
}
`, rInt)
}
