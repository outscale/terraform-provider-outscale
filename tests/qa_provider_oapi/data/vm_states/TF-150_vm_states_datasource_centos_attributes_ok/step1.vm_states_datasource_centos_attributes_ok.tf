resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF150"
}

resource "outscale_security_group" "security_group_TF150" {
  description         = "test-terraform-TF150"
  security_group_name = "terraform-sg-150"
}

resource "outscale_vm" "outscale_vm" {
   image_id                 = var.image_id
   vm_type                  = var.vm_type
   keypair_name             = outscale_keypair.my_keypair.keypair_name
   security_group_ids       = [outscale_security_group.security_group_TF150.security_group_id]
}

resource "outscale_vm" "outscale_vm2" {
    image_id                = var.image_id
    vm_type                 = var.vm_type
    keypair_name            = outscale_keypair.my_keypair.keypair_name
    security_group_ids      = [outscale_security_group.security_group_TF150.security_group_id]
}

data "outscale_vm_states" "outscale_vm_states" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id, outscale_vm.outscale_vm2.vm_id]
     }
}
