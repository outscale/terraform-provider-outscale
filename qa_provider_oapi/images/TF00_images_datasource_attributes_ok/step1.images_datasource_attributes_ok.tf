resource "outscale_vm" "outscale_vm" {
    image_id           = "ami-be23e98b"
}

resource "outscale_image" "outscale_image1" {
    image_name = "test datasources 1"
    vm_id      = outscale_vm.outscale_vm.vm_id
    no_reboot  = "true"
} 

resource "outscale_image" "outscale_image2" {
    image_name = "test datasources 2"
    vm_id      = outscale_vm.outscale_vm.vm_id
    no_reboot  = "true"
} 


data "outscale_images" "outscale_images" {

    #executable_by = ["339215505907"]    # customer-tooling

    filter {
		#name   = "architectures"
		#values = ["x86_64"]
	   	name   = "image_ids"
	   	values = [outscale_image.outscale_image1.image_id,outscale_image.outscale_image2.image_id]
	}
}
