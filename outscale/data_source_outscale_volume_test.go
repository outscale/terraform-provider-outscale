package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_VolumeDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumeDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "volume_type", "standard"),
				),
			},
		},
	})
}

func TestAccOthers_VolumeDataSource_filterByTags(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumeDataSourceConfigFilterByTags(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "volume_type", "standard"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVolumeDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Volume data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Volume data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleOAPIVolumeDataSourceConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 10

			tags {
				key   = "Name"
				value = "External Volume"
			}
		}

		data "outscale_volume" "ebs_volume" {
			filter {
				name   = "volume_ids"
				values = [outscale_volume.example.id]
			}
		}
	`, region)
}

func testAccCheckOutscaleOAPIVolumeDataSourceConfigFilterByTags(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 10

			tags {
				key   = "Name"
				value = "volume-io1-2"
			}
		}

		data "outscale_volume" "ebs_volume" {
			filter {
				name   = "tags"
				values = ["Name=volume-io1-2"]
			}

			filter {
				name   = "volume_ids"
				values = [outscale_volume.example.id]
			}

		}
	`, region)
}
