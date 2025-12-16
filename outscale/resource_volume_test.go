package outscale

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_Volume_Basic(t *testing.T) {
	t.Parallel()

	resourceName := "outscale_volume.accvolume"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_Volume_UpdateSize(t *testing.T) {
	t.Parallel()
	region := utils.GetRegion()

	resourceName := "outscale_volume.accvolume"
	var volumeID string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleVolumeConfig(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttrWith(resourceName, "volume_id", func(value string) error {
						volumeID = value
						return nil
					}),
				),
			},
			{
				Config: testOutscaleVolumeConfigUpdate(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "size", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttrWith(resourceName, "volume_id", func(value string) error {
						if value != volumeID {
							return fmt.Errorf("volume_id changed from %s to %s, resource was replaced instead of updated", volumeID, value)
						}
						return nil
					}),
				),
			},
		},
	})
}

func TestAccOthers_Volume_IO1Type(t *testing.T) {
	t.Parallel()

	region := utils.GetRegion()
	resourceName := "outscale_volume.test-io1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
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

func TestAccOthers_Volume_Type_Change(t *testing.T) {
	t.Parallel()

	region := utils.GetRegion()
	resourceName := "outscale_volume.test-type-change"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: test_VolumeTypeGP2Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "gp2"),
				),
			},
			{
				Config: test_VolumeTypeIO1Config(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "io1"),
					resource.TestCheckResourceAttr(resourceName, "iops", "100"),
				),
			},
			{
				Config: test_VolumeTypeSTDConfig(region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "volume_type", "standard"),
				),
			},
		},
	})
}

func TestAccOthers_Volume_Migration(t *testing.T) {
	region := utils.GetRegion()
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: FrameworkMigrationTestSteps("1.1.3",
			testAccOutscaleVolumeConfig(utils.GetRegion()),
			test_IO1VolumeTypeConfig(region),
		),
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

func test_VolumeTypeGP2Config(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-type-change" {
			subregion_name = "%sa"
			volume_type    = "gp2"
			size           = 10
		}
	`, region)
}

func test_VolumeTypeIO1Config(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-type-change" {
			subregion_name = "%sa"
			volume_type    = "io1"
			size           = 11
			iops           = 100
		}
	`, region)
}

func test_VolumeTypeSTDConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test-type-change" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 11
		}
	`, region)
}
