resource "outscale_vm" "outscale_vm" {
  # image_id                 = var.image_id
     image_id                 = "ami-be23e98b"
  #  vm_type                  = var.vm_type
     vm_type                  = "c4.large"
  #  keypair_name             = var.keypair_name
    keypair_name             = "integ_sut_keypair"
  #  security_group_ids       = [var.security_group_id]
}

resource "outscale_vm" "outscale_vm2" {
  # image_id                 = var.image_id
     image_id                 = "ami-be23e98b"
  #  vm_type                  = var.vm_type
     vm_type                  = "c4.large"
  #  keypair_name             = var.keypair_name
    keypair_name             = "integ_sut_keypair"
  #  security_group_ids       = [var.security_group_id]
}

data "outscale_vms_state" "outscale_vms_state" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.id, outscale_vm.outscale_vm2.id]
     }
}
