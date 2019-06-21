package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVolumeDataSource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volume.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "size", "10"),
					resource.TestCheckResourceAttr("data.outscale_volume.ebs_volume", "volume_type", "io1"),
					testAccCheckState("data.outscale_volume.ebs_volume"),
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

const testAccCheckOutscaleOAPIVolumeDataSourceConfig = `
resource "outscale_volume" "example" {
    subregion_name = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tags {
		key = "Name" 
		value = "External Volume"
	}
}

data "outscale_volume" "ebs_volume" {
    filter {
		name = "volume_ids"
		values = ["${outscale_volume.example.id}"]
    }
}
`
