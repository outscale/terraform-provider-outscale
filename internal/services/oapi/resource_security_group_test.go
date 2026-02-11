package oapi_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

var uniqueSGErr = regexp.MustCompile("The Security Group is the unique Security Group")

func TestAccNet_WithSecurityGroup(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "outscale_security_group.web"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "security_group_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroupWithoutName(t *testing.T) {
	resourceName := "outscale_security_group.noname"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSecurityGroupWithoutNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
		},
	})
}

func TestAccOthers_SecurityGroup_VM_CleanUp(t *testing.T) {
	resourceName := "outscale_security_group.group1"
	resourceName2 := "outscale_security_group.group2"
	sgName1 := acctest.RandomWithPrefix("testacc-sg")
	sgName2 := acctest.RandomWithPrefix("testacc-sg")
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: securityGroupWithVMConfigStep1(omi, testAccVmType, testAccKeypair, sgName1, sgName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Changing the SG name to unlink it from the VM and recreate it
				Config: securityGroupWithVMConfigStep2(omi, testAccVmType, testAccKeypair, sgName1, sgName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Removing one SG from the VM
				Config: securityGroupWithVMConfigStep3(omi, testAccVmType, testAccKeypair, sgName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
			{
				// Removing the default SG fails
				Config:      securityGroupWithVMConfigStep4(omi, testAccVmType, testAccKeypair),
				ExpectError: uniqueSGErr,
			},
		},
	})
}

func TestAccOthers_SecurityGroup_NIC_CleanUp(t *testing.T) {
	resourceName := "outscale_security_group.group1"
	resourceName2 := "outscale_security_group.group2"
	sgName1 := acctest.RandomWithPrefix("testacc-sg")
	sgName2 := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: securityGroupWithNICConfigStep1(sgName1, sgName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Changing the SG name to unlink it from the NIC and recreate it
				Config: securityGroupWithNICConfigStep2(sgName1, sgName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Removing one SG from the NIC
				Config: securityGroupWithNICConfigStep3(sgName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
			{
				// Removing the default SG fails
				Config:      securityGroupWithNICConfigStep4(),
				ExpectError: uniqueSGErr,
			},
		},
	})
}

func TestAccOthers_SecurityGroup_LBU_CleanUp(t *testing.T) {
	resourceName := "outscale_security_group.group1"
	resourceName2 := "outscale_security_group.group2"
	sgName1 := acctest.RandomWithPrefix("testacc-sg")
	sgName2 := acctest.RandomWithPrefix("testacc-sg")
	lbuName := acctest.RandomWithPrefix("testacc-lbu")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: securityGroupWithLBUConfigStep1(sgName1, sgName2, lbuName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Changing the SG name to unlink it from the LBU and recreate it
				Config: securityGroupWithLBUConfigStep2(sgName1, sgName2, lbuName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
					resource.TestCheckResourceAttrSet(resourceName2, "security_group_name"),
				),
			},
			{
				// Removing one SG from the LBU
				Config: securityGroupWithLBUConfigStep3(sgName1, lbuName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "security_group_name"),
				),
			},
			{
				// Removing the default LBU fails
				Config:      securityGroupWithLBUConfigStep4(lbuName),
				ExpectError: uniqueSGErr,
			},
		},
	})
}

func TestAccNet_WithSecurityGroup_Migration(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.2.1", testAccOutscaleSecurityGroupConfig(rInt), testAccOutscaleSecurityGroupWithoutNameConfig()),
	})
}

func testAccOutscaleSecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "web" {
			security_group_name = "terraform_test_%d"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = outscale_net.net.id
		}
	`, rInt)
}

func securityGroupWithLBUConfigStep1(sgName, sgName2, lbuName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "lbu security group"
			security_group_name = "%[1]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group2" {
			description         = "lbu security group"
			security_group_name = "%[2]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_load_balancer" "load_balancer02" {
			load_balancer_name = "%[3]s"
			listeners {
				backend_port           = 80
				backend_protocol       = "TCP"
				load_balancer_protocol = "TCP"
				load_balancer_port     = 80
			}
			subnets            = [outscale_subnet.subnet01.subnet_id]
			security_groups    = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
			load_balancer_type = "internal"
		}
	`, sgName, sgName2, lbuName)
}

func securityGroupWithLBUConfigStep2(sgName, sgName2, lbuName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "lbu security group"
			security_group_name = "%[1]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group2" {
			description         = "lbu security group"
			security_group_name = "%[2]s_recreated"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_load_balancer" "load_balancer02" {
			load_balancer_name = "%[3]s"
			listeners {
				backend_port           = 80
				backend_protocol       = "TCP"
				load_balancer_protocol = "TCP"
				load_balancer_port     = 80
			}
			subnets            = [outscale_subnet.subnet01.subnet_id]
			security_groups    = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
			load_balancer_type = "internal"
		}
	`, sgName, sgName2, lbuName)
}

func securityGroupWithLBUConfigStep3(sgName, lbuName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "lbu security group"
			security_group_name = "%[1]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_load_balancer" "load_balancer02" {
			load_balancer_name = "%[2]s"
			listeners {
				backend_port           = 80
				backend_protocol       = "TCP"
				load_balancer_protocol = "TCP"
				load_balancer_port     = 80
			}
			subnets            = [outscale_subnet.subnet01.subnet_id]
			security_groups    = [outscale_security_group.group1.security_group_id]
			load_balancer_type = "internal"
		}
	`, sgName, lbuName)
}

func securityGroupWithLBUConfigStep4(lbuName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_load_balancer" "load_balancer02" {
			load_balancer_name = "%[1]s"
			listeners {
				backend_port           = 80
				backend_protocol       = "TCP"
				load_balancer_protocol = "TCP"
				load_balancer_port     = 80
			}
			subnets            = [outscale_subnet.subnet01.subnet_id]
			security_groups    = []
			load_balancer_type = "internal"
		}
	`, lbuName)
}

func securityGroupWithNICConfigStep1(sgName, sgName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "nic security group"
			security_group_name = "%[1]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group2" {
			description         = "nic security group"
			security_group_name = "%[2]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_nic" "nic01" {
		    subnet_id = outscale_subnet.subnet01.subnet_id
		    security_group_ids = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
		}
	`, sgName, sgName2)
}

func securityGroupWithNICConfigStep2(sgName, sgName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "nic security group"
			security_group_name = "%[1]s_recreated"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group2" {
			description         = "nic security group"
			security_group_name = "%[2]s"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_nic" "nic01" {
		    subnet_id = outscale_subnet.subnet01.subnet_id
		    security_group_ids = [outscale_security_group.group1.security_group_id, outscale_security_group.group2.security_group_id]
		}
	`, sgName, sgName2)
}

func securityGroupWithNICConfigStep3(sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_security_group" "group1" {
			description         = "nic security group"
			security_group_name = "%[1]s_recreated"
			net_id              = outscale_net.net01.net_id
		}

		resource "outscale_nic" "nic01" {
		    subnet_id = outscale_subnet.subnet01.subnet_id
		    security_group_ids = [outscale_security_group.group1.security_group_id]
		}
	`, sgName)
}

func securityGroupWithNICConfigStep4() string {
	return `
		resource "outscale_net" "net01" {
    		ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "subnet01" {
		    subregion_name = "eu-west-2a"
		    ip_range       = "10.0.0.0/18"
		    net_id         = outscale_net.net01.net_id
		}

		resource "outscale_nic" "nic01" {
		    subnet_id = outscale_subnet.subnet01.subnet_id
		    security_group_ids = []
		}
	`
}

func securityGroupWithVMConfigStep1(omi, vmType, kpName, sgName, sgName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "group1" {
			description         = "vm security group"
			security_group_name = "%[4]s"
		}

		resource "outscale_security_group" "group2" {
			description         = "vm security group"
			security_group_name = "%[5]s"
		}

		resource "outscale_vm" "vm01" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[3]s"
			security_group_ids       = [outscale_security_group.group1.security_group_id,outscale_security_group.group2.security_group_id]
			placement_subregion_name = "eu-west-2a"
		}
	`, omi, vmType, kpName, sgName, sgName2)
}

func securityGroupWithVMConfigStep2(omi, vmType, kpName, sgName, sgName2 string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "group1" {
			description         = "vm security group"
			security_group_name = "%[4]s"
		}

		resource "outscale_security_group" "group2" {
			description         = "vm security group"
			security_group_name = "%[5]s_recreated"
		}

		resource "outscale_vm" "vm01" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[3]s"
			security_group_ids       = [outscale_security_group.group1.security_group_id,outscale_security_group.group2.security_group_id]
			placement_subregion_name = "eu-west-2a"
		}
	`, omi, vmType, kpName, sgName, sgName2)
}

func securityGroupWithVMConfigStep3(omi, vmType, kpName, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "group1" {
			description         = "vm security group"
			security_group_name = "%[4]s"
		}

		resource "outscale_vm" "vm01" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[3]s"
			security_group_ids       = [outscale_security_group.group1.security_group_id]
			placement_subregion_name = "eu-west-2a"
		}
	`, omi, vmType, kpName, sgName)
}

func securityGroupWithVMConfigStep4(omi, vmType, kpName string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "vm01" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[3]s"
			security_group_ids       = []
			placement_subregion_name = "eu-west-2a"
		}
	`, omi, vmType, kpName)
}

func testAccOutscaleSecurityGroupWithoutNameConfig() string {
	return `
		resource "outscale_security_group" "noname" {
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test-no-name"
			}
		}
	`
}
