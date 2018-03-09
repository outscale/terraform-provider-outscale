package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleEbsVolumesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleEbsVolumesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsEbsVolumeDataSourceID("data.outscale_volumes.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "size", "40"),
				),
			},
		},
	})
}

const testAccCheckOutscaleEbsVolumesDataSourceConfig = `
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tags {
        Name = "External Volume"
    }
}
data "outscale_volumes" "ebs_volume" {
    filter {
	name = "volume-type"
	values = ["${outscale_volume.example.volume_type}"]
    }
}
`
