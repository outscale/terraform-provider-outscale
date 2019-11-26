resource "outscale_image" "outscale_image" {
    image_name = "terraform test for image attributes"
    vm_id      = var.vm_id
    no_reboot  = "true"
}
