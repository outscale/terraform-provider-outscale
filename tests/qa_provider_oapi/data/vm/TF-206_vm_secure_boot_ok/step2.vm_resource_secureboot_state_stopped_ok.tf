resource "outscale_vm" "outscale_vm_TF206" {
  image_id            = var.image_id
  vm_type             = "tinav5.c3r3"
  deletion_protection = false
  boot_mode           = "uefi"
  secure_boot_action  = "enable"
  state               = "stopped"
  tags {
    key   = "name"
    value = "VM_TF206"
  }
}
