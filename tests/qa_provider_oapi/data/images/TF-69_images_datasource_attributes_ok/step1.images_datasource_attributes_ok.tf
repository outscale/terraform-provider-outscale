resource "outscale_vm" "outscale_vm" {
    image_id           = var.image_id
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
