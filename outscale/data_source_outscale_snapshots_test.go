package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_SnapshotsDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISnapshotsDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_snapshots.outscale_snapshots", "snapshots.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISnapshotsDataSourceConfig(region string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "example" {
		subregion_name = "%sa"
		size           = 1
	}

	resource "outscale_snapshot" "snapshot" {
		volume_id = outscale_volume.example.id
	}

	data "outscale_snapshots" "outscale_snapshots" {
		snapshot_id = [outscale_snapshot.snapshot.id]
	}
`, region)
}
