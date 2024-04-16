package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_VolumesDataSource_multipleFilters(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumeDataSourceConfigWithMultipleFilters(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volumes.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.size", "1"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.volume_type", "standard"),
				),
			},
		},
	})
}

func TestAccOthers_VolumeDataSource_multipleVIdsFilters(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumesDataSourceConfigWithMultipleVolumeIDsFilter(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volumes.outscale_volumes"),
					resource.TestCheckResourceAttr("data.outscale_volumes.outscale_volumes", "volumes.0.size", "40"),
				),
			},
		},
	})
}

func TestAccVM_withVolumesDataSource(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVolumesDataSourceConfigWithVM(utils.GetRegion(), omi, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVolumeDataSourceID("data.outscale_volumes.outscale_volumes"),
					// resource.TestCheckResourceAttr("data.outscale_volumes.outscale_volumes", "volumes.0.size", "1"),
					// resource.TestCheckResourceAttr("data.outscale_volumes.outscale_volumes", "volumes.0.volume_type", "standard"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIVolumeDataSourceConfigWithMultipleFilters(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "external" {
			subregion_name = "%sa"
			volume_type    = "standard"
			size           = 1

			tags {
				key   = "Name"
				value = "tf-acc-test-ebs-volume-test"
			}
		}

		data "outscale_volumes" "ebs_volume" {
			filter {
				name   = "volume_sizes"
				values = [outscale_volume.external.size]
			}

			filter {
				name   = "volume_types"
				values = [outscale_volume.external.volume_type]
			}
		}
	`, region)
}

func testAccCheckOutscaleOAPIVolumesDataSourceConfigWithMultipleVolumeIDsFilter(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%[1]sa"
			size           = 40
		}

		resource "outscale_volume" "outscale_volume2" {
			subregion_name = "%[1]sa"
			size           = 40
		}

		data "outscale_volumes" "outscale_volumes" {
			filter {
				name   = "volume_ids"
				values = [outscale_volume.outscale_volume.volume_id, outscale_volume.outscale_volume2.volume_id]
			}
		}
	`, region)
}

func testAccCheckOutscaleOAPIVolumesDataSourceConfigWithVM(region, imageID, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%[1]sa"
			volume_type    = "standard"
			size           = 25
			tags {
				key   = "Name"
				value = "volume-standard-1"
			}
		}

		resource "outscale_volume" "outscale_volume2" {
			subregion_name = "%[1]sa"
			volume_type    = "standard"
			size           = 13
			tags {
				key   = "Name"
				value = "volume-standard-2"
			}
		}

		resource "outscale_volume" "outscale_volume3" {
			subregion_name = "%[1]sa"
			size           = 40
			volume_type    = "gp2"
			tags {
				key   = "type"
				value = "gp2"
			}
		}

		resource "outscale_security_group" "sg_volumes" {
			security_group_name = "sg_vols"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "outscale_vm" {
			image_id           = "%[2]s"
			keypair_name       = "%[3]s"
			security_group_ids = [outscale_security_group.sg_volumes.security_group_id]
			vm_type            = "tinav4.c2r2p2"
		}

		resource "outscale_volumes_link" "outscale_volumes_link" {
			device_name = "/dev/xvdc"
			volume_id   = outscale_volume.outscale_volume.id
			vm_id       = outscale_vm.outscale_vm.id
		}

		resource "outscale_volumes_link" "outscale_volumes_link_2" {
			device_name = "/dev/xvdd"
			volume_id   = outscale_volume.outscale_volume2.id
			vm_id       = outscale_vm.outscale_vm.id
		}

		resource "outscale_volumes_link" "outscale_volumes_link_3" {
			device_name = "/dev/xvde"
			volume_id   = outscale_volume.outscale_volume3.id
			vm_id       = outscale_vm.outscale_vm.id
		}

		data "outscale_volumes" "outscale_volumes" {
			filter {
				name   = "link_volume_vm_ids"
				values = [outscale_vm.outscale_vm.vm_id]
			}
		}
	`, region, imageID, keypair)
}
