package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleVolumeDataSource_multipleFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumeDataSourceConfigWithMultipleFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volumes.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volume_set.0.size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volume_set.0.volume_type", "gp2"),
				),
			},
		},
	})
}

const testAccCheckOutscaleVolumeDataSourceConfigWithMultipleFilters = `
resource "outscale_volume" "external1" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 10
    tags {
        Name = "External Volume 1"
    }
}
data "outscale_volumes" "ebs_volume" {
    filter {
	name = "size"
	values = ["${outscale_volume.external1.size}"]
    }
    filter {
	name = "volume-type"
	values = ["${outscale_volume.external1.volume_type}"]
    }
}
`
