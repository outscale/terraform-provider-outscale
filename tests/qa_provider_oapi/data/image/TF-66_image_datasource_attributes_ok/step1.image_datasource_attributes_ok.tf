resource "outscale_vm" "my-vm" {
     image_id = var.image_id
}

resource "outscale_image" "outscale_image" {
    image_name = "terraform-image-1"
    vm_id      = outscale_vm.my-vm.vm_id
    no_reboot  = "true"
    tags {
      key = "Key-tags"
      value = "value-tags"
     }
tags {
      key = "Key-2"
      value = "value-tags-2"
     }
}


data "outscale_image" "outscale_image" {
    filter {
        name   = "image_ids"
        values = [outscale_image.outscale_image.image_id]
    }
}
