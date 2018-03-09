package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVolume_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "attachment_set.#", "0"),
				),
			},
		},
	})
}
func TestAccOutscaleOAPIVolume_NoIops(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfigWithNoIops,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.iops_test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolume_withTags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.tags_test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.tags_test", &v),
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

		var err error
		var response *fcu.DescribeVolumesOutput

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = conn.VM.DescribeVolumes(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		fmt.Printf("[DEBUG] Error Test Exists: %s", err)
		fmt.Printf("[DEBUG] Volume Exists: %v ", *response)
		if err == nil {
			if response.Volumes != nil && len(response.Volumes) > 0 {
				*v = *response.Volumes[0]
				return nil
			}
		}

		return fmt.Errorf("Error finding Outscale volume %s", rs.Primary.ID)
	}
}

const testAccOutscaleOAPIVolumeConfig = `
resource "outscale_volume" "test" {
  sub_region_name = "eu-west-2a"
  type = "gp2"
  size = 1
  tag = {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccOutscaleOAPIVolumeConfigWithTags = `
resource "outscale_volume" "tags_test" {
  sub_region_name = "eu-west-2a"
  size = 1
  tag = {
    Name = "TerraformTest"
  }
}
`

const testAccOutscaleOAPIVolumeConfigWithNoIops = `
resource "outscale_volume" "iops_test" {
  sub_region_name = "eu-west-2a"
  size = 10
  type = "gp2"
  iops = 0
  tag = {
    Name = "TerraformTest"
  }
}
`
