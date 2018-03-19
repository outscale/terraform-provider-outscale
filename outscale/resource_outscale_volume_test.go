package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVolume_basic(t *testing.T) {
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleVolume_updateSize(t *testing.T) {
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "1"),
				),
			},
			{
				Config: testOutscaleVolumeConfigUpdateSize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "10"),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string, v *fcu.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		request := &fcu.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(rs.Primary.ID)},
		}

		var response *fcu.DescribeVolumesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.VM.DescribeVolumes(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err == nil {
			if response.Volumes != nil && len(response.Volumes) > 0 {
				*v = *response.Volumes[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding EC2 volume %s", rs.Primary.ID)
	}
}

const testAccOutscaleVolumeConfig = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "gp2"
  size = 1
  tag {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testOutscaleVolumeConfigUpdateSize = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "gp2"
  size = 10
  tag {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`
