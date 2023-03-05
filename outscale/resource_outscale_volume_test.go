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

func TestAccOthers_Volume_basic(t *testing.T) {
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

func TestAccOthers_Volume_updateSize(t *testing.T) {
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

func TestAccOthers_Volume_io1Type(t *testing.T) {
	t.Parallel()
	if os.Getenv("IS_IO1_TEST_QUOTA") == "true" {
		var v oscgo.Volume
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)

			},
			IDRefreshName: "outscale_volume.test-io1",
			Providers:     testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: test_IO1VolumeTypeConfig(utils.GetRegion()),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckOAPIVolumeExists("outscale_volume.test-io1", &v),
					),
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func TestAccOthers_GP2_Volume_Type(t *testing.T) {
	t.Parallel()
	var v oscgo.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		IDRefreshName: "outscale_volume.test-gp2",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: test_GP2VolumeTypeConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVolumeExists("outscale_volume.test-gp2", &v),
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

func test_IO1VolumeTypeConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-io1" {
			subregion_name = "%sa"
			volume_type    = "gp2"
			size           = 10
			iops           = 100
		}
	`, region)
}

func test_GP2VolumeTypeConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-gp2" {
			subregion_name = "%sa"
			volume_type    = "gp2"
			size           = 10
		}
	`, region)
}
