package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleEbsVolumeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEbsVolumeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsEbsVolumeDataSourceID("data.outscale_volume.outscale_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.outscale_volume", "size", "40"),
				),
			},
		},
	})
}

func TestAccOutscaleEbsVolumeDataSource_multipleFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAwsEbsVolumeDataSourceConfigWithMultipleFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsEbsVolumeDataSourceID("data.outscale_volume.outscale_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.outscale_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.outscale_volume", "volume_type", "gp2"),
				),
			},
		},
	})
}

func testAccCheckAwsEbsVolumeDataSourceID(n string) resource.TestCheckFunc {
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

const testAccCheckAwsEbsVolumeDataSourceConfig = `
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tags = {
        Name = "External Volume"
    }
}
data "outscale_volume" "outscale_volume" {
    filter {
	name = "volume-type"
	values = ["${outscale_volume.example.volume_type}"]
    }
}
`

const testAccCheckAwsEbsVolumeDataSourceConfigWithMultipleFilters = `
resource "outscale_volume" "external1" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 10
    tags = {
        Name = "External Volume 1"
    }
}
data "outscale_volume" "outscale_volume" {
    filter = {
	name = "tag:Name"
	values = ["External Volume 1"]
    }
    filter = {
	name = "size"
	values = ["${outscale_volume.external1.size}"]
    }
    filter = {
	name = "volume-type"
	values = ["${outscale_volume.external1.type}"]
    }
}
`
