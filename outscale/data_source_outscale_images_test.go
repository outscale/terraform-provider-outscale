package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIImagesDataSource_Instance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIImagesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImagesDataSourceID("data.outscale_images.nat_ami"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIImagesDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find AMI data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AMI data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPIImagesDataSourceConfig = `
	data "outscale_images" "nat_ami" {
		filter {
			name = "architectures"
			values = ["x86_64"]
		}
	}
`
