package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_SnapshotsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSnapshotsDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_snapshots.outscale_snapshots", "snapshots.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSnapshotsDataSourceConfig(region string) string {
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
