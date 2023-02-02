package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAcc_VM_DataSource(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	vmType := "tinav4.c2r2p2"
	dataSourceVmName := "data.outscale_vm.vm"
	dataSourcesVmName := "data.outscale_vms.vms"
	dataSourceVmStateName := "data.outscale_vm_state.state"
	dataSourcesVmStateName := "data.outscale_vm_states.state"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_VM_DataSource_Config(omi, vmType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceVmName, "image_id", omi),
					resource.TestCheckResourceAttr(dataSourceVmName, "vm_type", vmType),
					resource.TestCheckResourceAttr(dataSourceVmName, "tags.#", "1"),

					resource.TestCheckResourceAttr(dataSourcesVmName, "vms.0.image_id", omi),
					resource.TestCheckResourceAttr(dataSourcesVmName, "vms.0.vm_type", vmType),

					resource.TestCheckResourceAttrSet(dataSourceVmStateName, "vm_id"),

					resource.TestCheckResourceAttr(dataSourcesVmStateName, "vm_states.#", "1"),
				),
			},
		},
	})
}

func testAcc_VM_DataSource_Config(omi, vmType string) string {
	return fmt.Sprintf(`
	resource "outscale_vm" "basic" {
		image_id			= "%s"
		vm_type				= "%s"
		keypair_name	= "terraform-basic"

		tags {
			key   = "name"
			value = "test acc"
		}
	}

	data "outscale_vms" "vms" {
		filter {
			name   = "vm_ids"
			values = [outscale_vm.basic.id]
		}
	}
    
	data "outscale_vm" "vm" {
		filter {
			name   = "vm_ids"
			values = [outscale_vm.basic.id]
		}
	}

	data "outscale_vm_state" "state" {
		filter {
			name   = "vm_ids"
			values = [outscale_vm.basic.id]
		}
	}

	data "outscale_vm_states" "state" {
		filter {
			name   = "vm_ids"
			values = [outscale_vm.basic.id]
		}
	}
	`, omi, vmType)
}
