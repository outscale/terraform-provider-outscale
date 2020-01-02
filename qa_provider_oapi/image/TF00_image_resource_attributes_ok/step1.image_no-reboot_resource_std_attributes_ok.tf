resource "outscale_image" "outscale_image" {
    image_name = "terraform test for image attributes"
    vm_id      = var.vm_id
    no_reboot  = "true"
}

### Test Copy Image

resource "outscale_image" "outscale_image_2" {
description = "Test-copy-image"
image_name = "test-copy-image"
source_image_id= var.image_id
source_region_name= var.region
}
