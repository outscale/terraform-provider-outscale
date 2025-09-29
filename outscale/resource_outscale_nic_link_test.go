package outscale

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_withNicLink_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	rInt := acctest.RandInt()
	resourceName := "outscale_nic_link.outscale_nic_link"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            resourceName,
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNicLinkConfigBasic(rInt, omi, utils.TestAccVmType, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "vm_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_id"),
				),
			},
		},
	})
}

func TestAccNet_ImportNicLink_Basic(t *testing.T) {
	resourceName := "outscale_nic_link.outscale_nic_link"
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNicLinkConfigBasic(rInt, omi, utils.TestAccVmType, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "device_number", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "vm_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleNicLinkStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func testAccCheckOutscaleNicLinkStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		log.Printf("LOG_ : %#+v\n", rs.Primary.Attributes["nic_id"])
		return rs.Primary.Attributes["nic_id"], nil
	}
}

func testAccOutscaleNicLinkConfigBasic(sg int, omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key   = "Name"
				value = "testacc-nic-link"
			}
		}

		resource "outscale_security_group" "security_group_nic" {
			security_group_name = "terraform_test_%d"
			description         = "Used in the terraform acceptance tests"
			net_id              = outscale_net.net.id

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = [outscale_security_group.security_group_nic.security_group_id]
			placement_subregion_name = "%[4]sa"
			subnet_id                = outscale_subnet.outscale_subnet.id
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[4]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.net.id
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids   = [outscale_security_group.security_group_nic.security_group_id]
		}

		resource "outscale_nic_link" "outscale_nic_link" {
			device_number = 1
			vm_id         = outscale_vm.vm.id
			nic_id        = outscale_nic.outscale_nic.id
		}
	`, sg, omi, vmType, region)
}
