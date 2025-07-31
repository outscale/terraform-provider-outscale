package outscale

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_Volume_basic(t *testing.T) {
	t.Parallel()

	resourceName := "outscale_volume.accvolume"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func TestAccOthers_Volume_updateSize(t *testing.T) {
	t.Parallel()
	region := utils.GetRegion()

	resourceName := "outscale_volume.accvolume"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				),
			},
			{
				Config: testOutscaleVolumeConfigUpdate(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
				),
			},
		},
	})
}

func TestAccOthers_Volume_io1Type(t *testing.T) {
	t.Parallel()
	region := utils.GetRegion()
	resourceName := "outscale_volume.test-io1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: test_IO1VolumeTypeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "io1"),
					resource.TestCheckResourceAttr(resourceName, "iops", "100"),
				),
			},
			{
				Config: test_IO1VolumeTypeConfigUpdate(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "io1"),
					resource.TestCheckResourceAttr(resourceName, "iops", "200"),
				),
			},
		},
	})
}

func TestAccOthers_GP2_Volume_Type(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_volume.test-gp2"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: test_GP2VolumeTypeConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "gp2"),
				),
			},
		},
	})
}

func testAccOutscaleVolumeConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "accvolume" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 1
		}
	`, region)
}

func testOutscaleVolumeConfigUpdate(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "accvolume" {
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
			volume_type    = "io1"
			size           = 10
			iops           = 100
		}
	`, region)
}
func test_IO1VolumeTypeConfigUpdate(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-io1" {
			subregion_name = "%sa"
			volume_type    = "io1"
			size           = 10
			iops           = 200
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
