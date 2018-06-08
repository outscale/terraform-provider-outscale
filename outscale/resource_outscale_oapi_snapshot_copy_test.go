package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshotCopy_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISnapshotCopyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotCopyExists("outscale_snapshot_copy.test"),
					resource.TestCheckNoResourceAttr("outscale_snapshot_copy.test", "snapshot_id"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISnapshotCopyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Snapshot copy id is set")
		}

		return nil
	}
}

func testAccOutscaleOAPISnapshotCopyConfig() string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" {
	sub_region_name = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
	description = "Snapshot Acceptance Test"
}

resource "outscale_snapshot_copy" "test" {
	source_region_name =  "eu-west-2b"
	source_snapshot_id = "${outscale_snapshot.test.id}"
}
`)
}
