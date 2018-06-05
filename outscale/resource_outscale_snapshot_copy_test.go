package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleSnapshotCopy_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleSnapshotCopyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotCopyExists("outscale_snapshot_copy.test"),
					resource.TestCheckNoResourceAttr("outscale_snapshot_copy.test", "snapshot_id"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSnapshotCopyExists(n string) resource.TestCheckFunc {
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

func testAccOutscaleSnapshotCopyConfig() string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" {
	availability_zone = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
	description = "Snapshot Acceptance Test"
}

resource "outscale_snapshot_copy" "test" {
	source_region =  "eu-west-2"
	source_snapshot_id = "${outscale_snapshot.test.id}"
}
`)
}
