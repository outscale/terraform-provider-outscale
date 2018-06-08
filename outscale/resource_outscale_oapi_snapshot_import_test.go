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
					testAccOutscaleOAPISnapshotCopyExists("outscale_snapshot_import"),
				),
			},
		},
	})
}

func testAccOutscaleOAPISnapshotCopyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccOutscaleOAPISnapshotCopyConfig() string {
	return fmt.Sprintf(`
resource "outscale_snapshot_import" "test" {
	osu_location = ""
snapshot_size = ""
}
`)
}
