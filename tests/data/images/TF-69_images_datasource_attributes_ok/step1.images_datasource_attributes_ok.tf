resource "outscale_security_group" "my_sgImgs" {
   description         = "test sg-group"
   security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "outscale_vm" {
   image_id           = var.image_id
   vm_type            = var.vm_type
   security_group_ids = [outscale_security_group.my_sgImgs.security_group_id]
}

resource "outscale_image" "outscale_image1" {
    image_name = "test-image-${random_string.suffix[0].result}"
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
