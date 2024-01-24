package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccVM_WithLBUAttachment_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.LoadBalancer
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	testCheckInstanceAttached := func(count int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if conf.BackendVmIds != nil {
				if len(*conf.BackendVmIds) != count {
					return fmt.Errorf("backend_vm_ids count does not match")
				}
				return nil
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAcc_ConfigLBUAttachmentBasic(omi, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(1),
				),
			},
			{
				Config: testAcc_ConfigLBUAttachmentAddUpdate(omi, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(2),
				),
			},
			{
				Config: testAcc_ConfigLBUAttachmentRemoveUpdate(omi, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
					testCheckInstanceAttached(1),
				),
			},
		},
	})
}

// add one attachment
func testAcc_ConfigLBUAttachmentBasic(omi, region string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
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
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_security_group" "sg_lb02" {
  security_group_name = "terraform_test_lb02"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo2" {
  image_id = "%[2]s"
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb02.security_group_id]
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = outscale_load_balancer.bar.id
  backend_vm_ids = [outscale_vm.foo1.id]
}
`, region, omi)
}

func testAcc_ConfigLBUAttachmentAddUpdate(omi, region string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
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
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_security_group" "sg_lb02" {
  security_group_name = "terraform_test_lb02"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo2" {
  image_id = "%[2]s"
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb02.security_group_id]
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = outscale_load_balancer.bar.id
  backend_vm_ids = [outscale_vm.foo1.id, outscale_vm.foo2.id]
}
`, region, omi)
}

func testAcc_ConfigLBUAttachmentRemoveUpdate(omi, region string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
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
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb1.security_group_id]
}

resource "outscale_security_group" "sg_lb02" {
  security_group_name = "terraform_test_lb02"
  description         = "Used in the terraform acceptance tests"
  tags {
	key   = "Name"
	value = "tf-acc-test"
	}
}

resource "outscale_vm" "foo2" {
  image_id = "%[2]s"
  vm_type = "tinav4.c1r1p1"
  security_group_ids   = [outscale_security_group.sg_lb02.security_group_id]
}

resource "outscale_load_balancer_vms" "foo1" {
  load_balancer_name      = outscale_load_balancer.bar.id
  backend_vm_ids = [outscale_vm.foo2.id]
}
`, region, omi)
}
