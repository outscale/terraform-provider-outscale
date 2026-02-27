package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_VolumeDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumeDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "volume_type", "standard"),
				),
			},
		},
	})
}

func TestAccOthers_VolumeDataSource_filterByTags(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumeDataSourceConfigFilterByTags(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "volume_type", "standard"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVolumeDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find volume data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("volume data source id not set")
		}
		return nil
	}
}

func testAccCheckOutscaleVolumeDataSourceConfig(region string) string {
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

func testAccCheckOutscaleVolumeDataSourceConfigFilterByTags(region string) string {
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
