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

func TestAccVM_PublicIPLink_VM_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_public_ip_link.ip_link"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIPLinkConfig(omi, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "link_public_ip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vm_id"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_PublicIPLink_NIC_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")
	resourceName := "outscale_public_ip_link.ip_link"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIPLinkNICConfig(omi, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "link_public_ip_id"),
					resource.TestCheckResourceAttrSet(resourceName, "nic_id"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccVM_PublicIPLink_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.4.0",
			testAccPublicIPLinkConfig(omi, sgName),
			testAccPublicIPLinkNICConfig(omi, sgName),
		),
	})
}

func testAccPublicIPLinkConfig(omi, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_link" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "vm_link" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = [outscale_security_group.sg_link.security_group_id]
			placement_subregion_name = "%[3]sa"

			lifecycle { ignore_changes = [state] }
		}

		resource "outscale_public_ip" "ip" {}

		resource "outscale_public_ip_link" "ip_link" {
			public_ip = outscale_public_ip.ip.public_ip
			vm_id     = outscale_vm.vm_link.id
		}
	`, omi, testAccVmType, utils.GetRegion(), testAccKeypair, sgName)
}

func testAccPublicIPLinkNICConfig(omi, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
		  ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		  subregion_name = "eu-west-2a"
		  ip_range       = "10.0.0.0/18"
		  net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "sg" {
		  	description         = "Terraform security group for nic with public IP link"
			security_group_name = "%[5]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_internet_service" "internet_service01" {}

		resource "outscale_internet_service_link" "internet_service_link01" {
		  internet_service_id = outscale_internet_service.internet_service01.internet_service_id
		  net_id              = outscale_net.net01.net_id
		}

		resource "outscale_nic" "nic01" {
		  subnet_id = outscale_subnet.subnet01.subnet_id
		  security_group_ids = [outscale_security_group.sg.security_group_id]
		}

		resource "outscale_public_ip" "ip" {}

		resource "outscale_public_ip_link" "ip_link" {
			public_ip = outscale_public_ip.ip.public_ip
			nic_id = outscale_nic.nic01.id
		}
	`, omi, testAccVmType, utils.GetRegion(), testAccKeypair, sgName)
}
