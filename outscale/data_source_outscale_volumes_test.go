package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOutscaleOAPIVolumesDataSource_multipleFilters(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumeDataSourceConfigWithMultipleFilters(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volumes.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.size", "1"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.volume_type", "gp2"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolumeDataSource_multipleVIdsFilters(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumesDataSourceConfigWithMultipleVolumeIDsFilter(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volumes.outscale_volumes"),
					resource.TestCheckResourceAttr("data.outscale_volumes.outscale_volumes", "volumes.0.size", "40"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVolumeDataSourceConfigWithMultipleFilters(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "external" {
			subregion_name = "%sa"
			volume_type    = "gp2"
			size           = 1
		
			tags {
				key   = "Name"
				value = "tf-acc-test-ebs-volume-test"
			}
		}

		data "outscale_volumes" "ebs_volume" {
			filter {
				name   = "volume_sizes"
				values = ["${outscale_volume.external.size}"]
			}

			filter {
				name   = "volume_types"
				values = ["${outscale_volume.external.volume_type}"]
			}
		}
	`, region)
}

func testAccCheckOutscaleOAPIVolumesDataSourceConfigWithMultipleVolumeIDsFilter(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%[1]sa"
			size           = 40
		}

		resource "outscale_volume" "outscale_volume2" {
			subregion_name = "%[1]sa"
			size           = 40
		}

		data "outscale_volumes" "outscale_volumes" {
			filter {
				name   = "volume_ids"
				values = ["${outscale_volume.outscale_volume.volume_id}", "${outscale_volume.outscale_volume2.volume_id}"]
			}
		}
	`, region)
}
