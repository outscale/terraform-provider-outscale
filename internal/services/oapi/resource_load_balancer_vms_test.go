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

func TestAccVM_LbuBackends_Basic(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	resourceName := "outscale_load_balancer_vms.backend_test"
	sgName := acctest.RandomWithPrefix("testacc-sg")
	lbName := acctest.RandomWithPrefix("testacc-lbu")
	region := utils.GetRegion()
	vmType := testAccVmType

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBUAttachmentConfig1(lbName, omi, region, vmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_vm_ids.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_vm_ids.#", "2"),
				),
			},
			{
				Config: testAccLBUAttachmentAddUpdate(lbName, omi, region, vmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_vm_ids.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_vm_ids.#", "1"),
				),
			},
			{
				Config: testAccLBUAttachmentBackendIps(lbName, omi, region, vmType, sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_ips.#"),
					resource.TestCheckResourceAttr(resourceName, "backend_ips.#", "2"),
				),
			},
		},
	})
}

func TestAccVM_LbuBackends_Migration(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	vmType := testAccVmType
	sgName := acctest.RandomWithPrefix("testacc-sg")
	lbName := acctest.RandomWithPrefix("testacc-lbu")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.1.3",
			testAccLBUAttachmentConfig1(lbName, omi, region, vmType, sgName),
			testAccLBUAttachmentAddUpdate(lbName, omi, region, vmType, sgName),
		),
	})
}

// add one attachment
func testAccLBUAttachmentConfig1(lbName, omi, region, vmType, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = ["%[2]sa"]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%[5]s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = "%[3]s"
  vm_type = "%[4]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1[0].vm_id, outscale_vm.foo1[1].vm_id]
}
`, lbName, region, omi, vmType, sgName)
}

func testAccLBUAttachmentAddUpdate(lbName, omi, region, vmType, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = ["%[2]sa"]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%[5]s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = "%[3]s"
  vm_type = "%[4]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1[0].vm_id]
}
`, lbName, region, omi, vmType, sgName)
}

func testAccLBUAttachmentBackendIps(lbName, omi, region, vmType, sgName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "%s"
	subregion_names = ["%[2]sa"]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "%[5]s"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  count = 2
  image_id = "%[3]s"
  vm_type = "%[4]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_ips = [outscale_vm.foo1[0].public_ip, outscale_vm.foo1[1].public_ip]
}
`, lbName, region, omi, vmType, sgName)
}
