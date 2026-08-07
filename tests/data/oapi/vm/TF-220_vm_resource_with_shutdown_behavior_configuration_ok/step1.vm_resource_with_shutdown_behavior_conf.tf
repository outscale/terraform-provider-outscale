resource "outscale_vm" "outscale_vm" {
  image_id             = var.image_id
  vm_type              = var.vm_type
  keypair_name         = var.keypair_name
  vm_initiated_shutdown_behavior = "stop"
}
