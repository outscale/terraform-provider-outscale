package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVolume_basic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
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
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
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
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		IDRefreshName: "outscale_volume.test-io",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testOutscaleOAPIVolumeConfigIO1Type(region),
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
			response, _, err = conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(request)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
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
			volume_type    = "gp2"
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
			volume_type    = "gp2"
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

func testAccOutscaleOAPIVolumeConfigUpdateTags(region, value string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "outscale_volume" {
	  volume_type = "gp2"
	  subregion_name = "%sa"
	  size = 10
	  tags {
		key = "name" 
		value = "%s"
	  }
	}`, region, value)
}
