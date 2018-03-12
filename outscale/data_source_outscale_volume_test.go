package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVolumeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "40"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVolumeDataSourceID(n string) resource.TestCheckFunc {
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

const testAccCheckOutscaleVolumeDataSourceConfig = `
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tags {
        Name = "External Volume"
    }
}
data "outscale_volume" "ebs_volume" {
    filter {
		name = "volume-id"
		values = ["${outscale_volume.example.id}"]
    }
}
`
