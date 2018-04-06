package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVMAttr_Basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: "outscale_vm.outscale_vm",
		CheckDestroy:  testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_vm_attributes.outscale_vm_attributes", "ebs_optimized", "false"),
				),
			},
		},
	})
}

func testAccCheckOutscaleVMAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return nil
	}
}

func testAccCheckOutscaleVMConfig_basic() string {
	return `
resource "outscale_vm" "outscale_vm" {

    count = 1

    image_id = "ami-880caa66"

    instance_type = "c4.large"

    #key_name = "integ_sut_keypair"

    #security_group = ["sg-c73d3b6b"]

		disable_api_termination = true
		
		#ebs_optimized = true

}

 

resource "outscale_vm_attributes" "outscale_vm_attributes" {

    instance_id = "${outscale_vm.outscale_vm.0.id}"

    attribute = "disableApiTermination"
		disable_api_termination = false
		
		#attribute = "instanceType"
		#instance_type = "t2.micro"
		
    #attribute = "ebsOptimized"
		#ebs_optimized = false
		
		#attribute = "blockDeviceMapping"
		#block_device_mapping {
		#	device_name = "/dev/sda1"
		#		ebs {
		#			delete_on_termination = true
		#		}
		#}

}`
}
