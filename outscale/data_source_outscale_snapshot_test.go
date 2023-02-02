package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Snapshot_DataSource(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")
	dataSourceName := "data.outscale_snapshot.snapshot"
	dataSourcesName := "data.outscale_snapshots.snapshot"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Snapshot_DataSource_Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourcesName, "snapshots.#"),

					resource.TestCheckResourceAttr(dataSourceName, "volume_size", "1"),
				),
			},
		},
	})
}

func testAcc_Snapshot_DataSource_Config(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%sa"
			size           = 1
		}

		resource "outscale_snapshot" "snapshot" {
			volume_id = outscale_volume.example.id
		}

		data "outscale_snapshot" "snapshot" {
			filter {
				name   = "volume_ids"
				values = [outscale_snapshot.snapshot.volume_id]
			}
			filter {
				name   = "volume_sizes"
				values = [outscale_snapshot.snapshot.volume_size]
			}
		}

		data "outscale_snapshots" "snapshot" {
			filter {
				name   = "volume_ids"
				values = [outscale_snapshot.snapshot.volume_id]
			}
			filter {
				name   = "volume_sizes"
				values = [outscale_snapshot.snapshot.volume_size]
			}
		}
	`, region)
}
