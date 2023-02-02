package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_Volume_DataSource(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")

	dataSourceName := "data.outscale_volume.volume"
	dataSourcesName := "data.outscale_volumes.volumes"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Volume_DataSource_Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "size", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "volume_type", "standard"),

					resource.TestCheckResourceAttr(dataSourcesName, "volumes.#", "2"),
				),
			},
		},
	})
}

func testAcc_Volume_DataSource_Config(region string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "vol" {
		subregion_name = "%[1]sa"
		volume_type    = "standard"
		size           = 2

		tags {
			key   = "Name"
			value = "volume-standard-1"
		}
	}

	resource "outscale_volume" "vol2" {
		subregion_name = "%[1]sa"
		size           = 4
		iops           = 100
		volume_type    = "io1"
		tags {
			key   = "type"
			value = "io1"
		}
	}

	data "outscale_volume" "volume" {
		filter {
			name   = "tag_values"
			values = ["volume-standard-1"]
		}
		filter {
			name   = "volume_ids"
			values = ["${outscale_volume.vol.id}"]
		}
		filter {
			name   = "volume_sizes"
			values = ["${outscale_volume.vol.size}"]
		}
	}

	data "outscale_volumes" "volumes" {
		filter {
			name   = "volume_ids"
			values = ["${outscale_volume.vol.volume_id}", "${outscale_volume.vol2.volume_id}"]
		}

		filter {
			name   = "volume_sizes"
			values = ["${outscale_volume.vol.size}", "${outscale_volume.vol2.size}"]
		}

		filter {
			name   = "volume_types"
			values = ["${outscale_volume.vol.volume_type}", "${outscale_volume.vol2.volume_type}"]
		}
	}
	`, region)
}
