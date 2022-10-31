package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccVM_WithLBUAttachment_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.LoadBalancer
	omi := os.Getenv("OUTSCALE_IMAGEID")
	resourceName := "outscale_load_balancer_vms.backend_test"
	rand := acctest.RandIntRange(0, 50)
	region := utils.GetRegion()
	vmType := utils.TestAccVmType

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckValues(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUAttachmentConfig1(rand, omi, region, vmType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.lbu_test", &conf),
				),
			},
			{
				Config: testAcc_ConfigLBUAttachmentAddUpdate(omi, region, vmType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_ips.#"),
				),
			},
			{
				Config: testAcc_ConfigLBUAttachmentRemoveUpdate(omi, region, vmType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backend_vm_ids.0"),
					resource.TestCheckResourceAttr(resourceName, "backend_ips.#", "0"),
				),
			},
		},
	})
}

// add one attachment
func testAccOutscaleLBUAttachmentConfig1(num int, omi, region, vmType string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
	load_balancer_name = "load-test-%d"
	subregion_names = ["%[2]sa"]
    listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "terraform_test_lb01"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  image_id = "%[3]s"
  vm_type = "%[4]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_vm" "foo2" {
  image_id = "%[3]s"
  vm_type = "%[4]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1.vm_id]
}
`, num, region, omi, vmType)
}

func testAcc_ConfigLBUAttachmentAddUpdate(omi, region, vmType string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
  load_balancer_name = "load-test12"
  subregion_names = ["%[1]sa"]
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "terraform_test_lb01"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  image_id = "%[2]s"
  vm_type = "%[3]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_vm" "foo2" {
  image_id = "%[2]s"
  vm_type = "%[3]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1.vm_id]
  backend_ips = [outscale_vm.foo2.public_ip]
}
`, region, omi, vmType)
}

func testAcc_ConfigLBUAttachmentRemoveUpdate(omi, region, vmType string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "lbu_test" {
  load_balancer_name = "load-test12"
  subregion_names = ["%sa"]
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }
}

resource "outscale_security_group" "sg_lb1" {
  security_group_name = "terraform_test_lb01"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo1" {
  image_id = "%[2]s"
  vm_type = "%[3]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_vm" "foo2" {
  image_id = "%[2]s"
  vm_type = "%[3]s"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_load_balancer_vms" "backend_test" {
  load_balancer_name      = outscale_load_balancer.lbu_test.load_balancer_name
  backend_vm_ids = [outscale_vm.foo1.vm_id]
}
`, region, omi, vmType)
}
