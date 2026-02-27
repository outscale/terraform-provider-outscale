package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_VolumesDataSource_multipleFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumeDataSourceConfigWithMultipleFilters(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volumes.ebs_volume"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.size", "1"),
					resource.TestCheckResourceAttr("data.outscale_volumes.ebs_volume", "volumes.0.volume_type", "standard"),
				),
			},
		},
	})
}

func TestAccOthers_VolumeDataSource_multipleVIdsFilters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumesDataSourceConfigWithMultipleVolumeIDsFilter(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volumes.outscale_volumes"),
					resource.TestCheckResourceAttr("data.outscale_volumes.outscale_volumes", "volumes.0.size", "40"),
				),
			},
		},
	})
}

func TestAccVM_withVolumesDataSource(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVolumesDataSourceConfigWithVM(utils.GetRegion(), omi, keypair, testAccVmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVolumeDataSourceID("data.outscale_volumes.outscale_volumes"),
					resource.TestCheckResourceAttrSet("data.outscale_volumes.outscale_volumes", "volumes.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVolumeDataSourceConfigWithMultipleFilters(region string) string {
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

func testAccCheckOutscaleVolumesDataSourceConfigWithMultipleVolumeIDsFilter(region string) string {
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

func testAccCheckOutscaleVolumesDataSourceConfigWithVM(region, imageID, keypair, vmType, sgName string) string {
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
			security_group_name = "%[5]s"
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
			vm_type            = "%[4]s"
		}

		resource "outscale_volume_link" "outscale_volume_link" {
			device_name = "/dev/xvdc"
			volume_id   = outscale_volume.outscale_volume.id
			vm_id       = outscale_vm.outscale_vm.id
		}

		resource "outscale_volume_link" "outscale_volume_link_2" {
			device_name = "/dev/xvdd"
			volume_id   = outscale_volume.outscale_volume2.id
			vm_id       = outscale_vm.outscale_vm.id
		}

		resource "outscale_volume_link" "outscale_volume_link_3" {
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
	`, region, imageID, keypair, vmType, sgName)
}
