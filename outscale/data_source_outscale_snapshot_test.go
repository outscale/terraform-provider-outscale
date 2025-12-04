package outscale

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_SnapshotDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSnapshotDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotDataSourceID("data.outscale_snapshot.snapshot_test"),
					resource.TestCheckResourceAttr("data.outscale_snapshot.snapshot_test", "volume_size", "1"),
				),
			},
		},
	})
}

func TestAccOthers_SnapshotDataSource_multipleFilters(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleSnapshotDataSourceConfigWithMultipleFilters(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleSnapshotDataSourceID("data.outscale_snapshot.snapshot_filters"),
					resource.TestCheckResourceAttr("data.outscale_snapshot.snapshot_filters", "volume_size", "10"),
				),
			},
		},
	})
}

func testAccCheckOutscaleSnapshotDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find snapshot data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleSnapshotDataSourceConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%sa"
			size           = 1
		}

		resource "outscale_snapshot" "snapshot_01" {
			volume_id = outscale_volume.example.id
		}

		data "outscale_snapshot" "snapshot_test" {
			snapshot_id = outscale_snapshot.snapshot_01.id
		}
	`, region)
}

func testAccCheckOutscaleSnapshotDataSourceConfigWithMultipleFilters(region string) string {
	creationDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`
		resource "outscale_volume" "external1" {
			subregion_name = "%sa"
			size           = 10
		}

		resource "outscale_snapshot" "snapshot_t2" {
			volume_id = outscale_volume.external1.id
		}

		data "outscale_snapshot" "snapshot_filters" {
			snapshot_id = outscale_snapshot.snapshot_t2.id

			filter {
				name   = "volume_sizes"
				values = ["10"]
			}
			filter {
				name   = "from_creation_date"
				values = ["%s"]
			}
		}
	`, region, creationDate)
}
