resource "outscale_security_group" "my_sgImgs" {
   description         = "test sg-group"
   security_group_name = "SG-TF69"
}

resource "outscale_vm" "outscale_vm" {
   image_id           = var.image_id
   vm_type            = var.vm_type
   security_group_ids = [outscale_security_group.my_sgImgs.security_group_id]
}

resource "outscale_image" "outscale_image1" {
    image_name = "TF-69-name"
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
