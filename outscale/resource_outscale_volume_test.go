package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIVolume_basic(t *testing.T) {
	t.Parallel()

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIVolumeConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test", &v),
				),
			},
			{
				ResourceName: "outscale_volume.test",
				ImportState:  true,
				//ImportStateVerify: true,
			},
		},
	})
}

func TestAccOutscaleOAPIVolume_updateSize(t *testing.T) {
	t.Parallel()
	region := utils.GetRegion()

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
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
	t.Parallel()

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		IDRefreshName: "outscale_volume.test-io",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOutscaleOAPIVolumeConfigIO1Type(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test-io", &v),
				),
			},
		},
	})
}

func testAccCheckOAPIVolumeExists(n string, v *oscgo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{rs.Primary.ID}},
		}

		var response oscgo.ReadVolumesResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			response = rp
			return nil
		})

		if err == nil {
			if response.Volumes != nil && len(response.GetVolumes()) > 0 {
				*v = response.GetVolumes()[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding volume %s", rs.Primary.ID)
	}
}

func testAccOutscaleOAPIVolumeConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 1
		
			tags {
				key   = "Name"
				value = "tf-acc-test-ebs-volume-test"
			}
		}
	`, region)
}

func testOutscaleOAPIVolumeConfigUpdateSize(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 10
		
			tags {
				key   = "Name"
				value = "tf-acc-test-ebs-volume-test"
			}
		}
	`, region)
}

func testOutscaleOAPIVolumeConfigIO1Type(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-io" {
			subregion_name = "%sa"
			volume_type    = "io1"
			size           = 10
			iops           = 100
		}
	`, region)
}
