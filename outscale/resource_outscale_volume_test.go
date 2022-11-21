package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVolume_basic(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
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

func TestAccVolume_updateSize(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "1"),
				),
			},
			{
				Config: testVolumeConfigUpdateSize(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "size", "10"),
				),
			},
		},
	})
}

func TestAccVolume_io1Type(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		IDRefreshName: "outscale_volume.test-io",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVolumeConfigIO1Type(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test-io", &v),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string, v *oscgo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*Client).OSCAPI

		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{rs.Primary.ID}},
		}

		var response oscgo.ReadVolumesResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
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

func testAccVolumeConfig(region string) string {
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

func testVolumeConfigUpdateSize(region string) string {
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

func testVolumeConfigIO1Type(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-io" {
			subregion_name = "%sa"
			volume_type    = "io1"
			size           = 10
			iops           = 100
		}
	`, region)
}
