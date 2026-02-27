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

func TestAccNet_WithInternetServiceLink_Basic(t *testing.T) {
	resourceName := "outscale_internet_service_link.outscale_internet_service_link"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_service_id"),
				),
			},
		},
	})
}

func TestAccVm_WithInternetServiceLink_Unlink_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	resourceName := "outscale_vm.outscale_vm_islink"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigWithInternetService(omi, testAccVmType, utils.GetRegion(), testAccKeypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", testAccVmType),
				),
			},
			{
				Config: testAccCheckOutscaleVMConfigWithInternetServiceRemoved(omi, testAccVmType, utils.GetRegion(), testAccKeypair, sgName),
			},
		},
	})
}

func TestAccNet_WithImportInternetServiceLink_Basic(t *testing.T) {
	resourceName := "outscale_internet_service_link.outscale_internet_service_link"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func testAccCheckOutscaleVMConfigWithInternetService(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags  {
			key   = "name"
			value = "Terraform_net"
		}
	}
	resource "outscale_subnet" "outscale_subnet" {
		net_id         = outscale_net.outscale_net.net_id
		ip_range       = "10.0.0.0/24"
		subregion_name = "%[3]sa"

		tags {
			key   = "name"
			value = "Terraform_subnet"
		}
	}

	resource "outscale_security_group" "outscale_sg" {
		description         = "sg for terraform tests"
		security_group_name = "%[5]s"
		net_id              = outscale_net.outscale_net.net_id
	}

	resource "outscale_internet_service" "outscale_internet_service" {
                depends_on = [outscale_net.outscale_net]
        }

	resource "outscale_route_table" "outscale_route_table" {
		net_id = outscale_net.outscale_net.net_id
		tags {
			key   = "name"
			value = "Terraform_RT"
		}
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		route_table_id  = outscale_route_table.outscale_route_table.route_table_id
		subnet_id       = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
		net_id              = outscale_net.outscale_net.net_id
	}

	resource "outscale_route" "outscale_route" {
		gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
		destination_ip_range = "0.0.0.0/0"
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}
	resource "outscale_vm" "outscale_vm_islink" {
		image_id           = "%[1]s"
		vm_type            = "%[2]s"
		keypair_name       = "%[4]s"
		security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		subnet_id          = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_public_ip_link" "outscale_public_ip_link" {
		vm_id     = outscale_vm.outscale_vm_islink.vm_id
		public_ip = outscale_public_ip.outscale_public_ip.public_ip
	}
	`, omi, vmType, region, keypair, sgName)
}

func testAccCheckOutscaleVMConfigWithInternetServiceRemoved(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags  {
			key   = "name"
			value = "Terraform_net"
		}
	}
	resource "outscale_subnet" "outscale_subnet" {
		net_id         = outscale_net.outscale_net.net_id
		ip_range       = "10.0.0.0/24"
		subregion_name = "%[3]sa"

		tags {
			key   = "name"
			value = "Terraform_subnet"
		}
	}

	resource "outscale_security_group" "outscale_sg" {
		description         = "sg for terraform tests"
		security_group_name = "%[5]s"
		net_id              = outscale_net.outscale_net.net_id
	}

	resource "outscale_internet_service" "outscale_internet_service" {
                depends_on = [outscale_net.outscale_net]
        }

	resource "outscale_route_table" "outscale_route_table" {
		net_id = outscale_net.outscale_net.net_id
		tags {
			key   = "name"
			value = "Terraform_RT"
		}
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		route_table_id  = outscale_route_table.outscale_route_table.route_table_id
		subnet_id       = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_route" "outscale_route" {
		gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
		destination_ip_range = "0.0.0.0/0"
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}
	resource "outscale_vm" "outscale_vm_islink" {
		image_id           = "%[1]s"
		vm_type            = "%[2]s"
		keypair_name       = "%[4]s"
		security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		subnet_id          = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_public_ip_link" "outscale_public_ip_link" {
		vm_id     = outscale_vm.outscale_vm_islink.vm_id
		public_ip = outscale_public_ip.outscale_public_ip.public_ip
	}
	`, omi, vmType, region, keypair, sgName)
}

func TestAccNet_WithInternetServiceLink_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.1.3", testAccOutscaleInternetServiceLinkConfig()),
	})
}

func testAccOutscaleInternetServiceLinkConfig() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-internet-service-link-rs"
			}
		}

		resource "outscale_internet_service" "outscale_internet_service" {}

		resource "outscale_internet_service_link" "outscale_internet_service_link" {
			net_id              = outscale_net.outscale_net.net_id
			internet_service_id = outscale_internet_service.outscale_internet_service.id
		}
	`
}
