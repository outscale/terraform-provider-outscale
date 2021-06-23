resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
}

resource "outscale_image" "outscale_image1" {
#    count      = 2
#    image_name = "test-datasources-${count.index}"
    image_name = "test-1"
    vm_id      = outscale_vm.outscale_vm.vm_id
    no_reboot  = "true"
    tags {
       key = "Key"
       value = "value-tags"
     }
    tags {
       key = "Key-2"
       value = "value-tags-2"
     }
} 


#data "outscale_images" "outscale_images" {
#
#    filter {
#	   	name   = "image_ids"
#	   	values = [outscale_image.outscale_image1[0].image_id,outscale_image.outscale_image1[1].image_id]
#	}
#}
