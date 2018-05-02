package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleSnapshotsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSnapshotsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_snapshots.outscale_snapshots", "snapshot_set.#", "1"),
				),
			},
		},
	})
}

const testAccCheckOutscaleSnapshotsDataSourceConfig = `
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tag {
        Name = "External Volume"
    }
}

resource "outscale_snapshot" "snapshot" {
    volume_id = "${outscale_volume.example.id}"
}

data "outscale_snapshots" "outscale_snapshots" {
    snapshot_id = ["${outscale_snapshot.snapshot.id}"]
}
`
