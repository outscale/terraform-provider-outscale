resource "outscale_vm" "outscale_vm" {
   image_id                 = var.image_id
   vm_type                  = var.vm_type
   keypair_name             = var.keypair_name
   security_group_ids       = [var.security_group_id]
}

resource "outscale_vm" "outscale_vm2" {
    image_id                  = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = var.keypair_name 
    security_group_ids       = [var.security_group_id]
}

data "outscale_vms_state" "outscale_vms_state" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id, outscale_vm.outscale_vm2.vm_id]
     }
}
