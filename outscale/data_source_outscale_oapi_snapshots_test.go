package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleOAPISnapshotsDataSource_basic(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISnapshotsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_snapshots.outscale_snapshots", "snapshot_set.#", "1"),
				),
			},
		},
	})
}

const testAccCheckOutscaleOAPISnapshotsDataSourceConfig = `
	resource "outscale_volume" "example" {
		subregion_name = "in-west-2a"
		size           = 1
	}

	resource "outscale_snapshot" "snapshot" {
		volume_id = "${outscale_volume.example.id}"
	}

	data "outscale_snapshots" "outscale_snapshots" {
		snapshot_id = ["${outscale_snapshot.snapshot.id}"]
	}
`
