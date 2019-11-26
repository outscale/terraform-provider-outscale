resource "outscale_vm" "outscale_vm" {
    image_id            = var.image_id
    vm_type             = var.vm_type
    keypair_name        = var.keypair_name
    security_group_ids  = [var.security_group_id]
    deletion_protection = true
}

resource "outscale_vm_attributes" "outscale_vm_attributes" {
    vm_id               = outscale_vm.outscale_vm.vm_id
    deletion_protection = false
} 
