resource "outscale_security_group" "my_sgImg1" {
   description = "test sg group"
}

resource "outscale_vm" "my-vm" {
    image_id           = var.image_id
    security_group_ids = [outscale_security_group.my_sgImg1.security_group_id]
}

resource "outscale_image" "outscale_image" {
    image_name = "terraform-image-2"
    vm_id      = outscale_vm.my-vm.vm_id
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

### Test Copy Image

resource "outscale_image" "outscale_image_2" {
description = "Test-copy-image"
image_name = "terraform-copy-image"
source_image_id= var.image_id
source_region_name= var.region
}
