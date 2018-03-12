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

func TestAccOutscaleOAPIVolume_basic(t *testing.T) {
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolume_updateSize(t *testing.T) {
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "1"),
				),
			},
			{
				Config: testOutscaleOAPIVolumeConfigUpdateSize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "10"),
				),
			},
		},
	})
}

func testAccCheckOAPIVolumeExists(n string, v *fcu.Volume) resource.TestCheckFunc {
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

const testAccOutscaleOAPIVolumeConfig = `
resource "outscale_volume" "test" {
  sub_region_name = "eu-west-2a"
  type = "gp2"
  size = 1
  tag {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testOutscaleOAPIVolumeConfigUpdateSize = `
resource "outscale_volume" "test" {
  sub_region_name = "eu-west-2a"
  type = "gp2"
  size = 10
  tag {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`
