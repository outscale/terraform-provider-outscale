package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func TestAccOutscaleOAPIVolume_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")
	region := os.Getenv("OUTSCALE_REGION")

	if region == "" {
		region = "dv-west-1"
	}

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	var v oapi.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolume_updateSize(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")
	region := os.Getenv("OUTSCALE_REGION")

	if region == "" {
		region = "dv-west-1"
	}

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	var v oapi.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "1"),
				),
			},
			{
				Config: testOutscaleOAPIVolumeConfigUpdateSize(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "10"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVolume_io1Type(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	var v oapi.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test-io",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOutscaleOAPIVolumeConfigIO1Type,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test-io", &v),
				),
			},
		},
	})
}

func testAccCheckOAPIVolumeExists(n string, v *oapi.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		request := &oapi.ReadVolumesRequest{
			Filters: oapi.FiltersVolume{VolumeIds: []string{rs.Primary.ID}},
		}

		var response *oapi.ReadVolumesResponse
		var resp *oapi.POST_ReadVolumesResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadVolumes(*request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		response = resp.OK

		if err == nil {
			if response.Volumes != nil && len(response.Volumes) > 0 {
				*v = response.Volumes[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding EC2 volume %s", rs.Primary.ID)
	}
}

func testAccOutscaleOAPIVolumeConfig(region string) string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" {
  sub_region_name = "%sa"
  volume_type = "gp2"
  size = 1
  tags {
	key = "Name" 
	value = "tf-acc-test-ebs-volume-test"
  }
}
`, region)
}

func testOutscaleOAPIVolumeConfigUpdateSize(region string) string {
	return fmt.Sprintf(`
resource "outscale_volume" "test" {
  sub_region_name = "%sa"
  volume_type = "gp2"
  size = 10
  tags {
	key = "Name" 
	value = "tf-acc-test-ebs-volume-test"
  }
}
`, region)
}

const testOutscaleOAPIVolumeConfigIO1Type = `
resource "outscale_volume" "test-io" {
	sub_region_name = "dv-west-1a"
	size = 10
	iops = 5
	volume_type = "io1"
}

`
