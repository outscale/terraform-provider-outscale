resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_security_group" "security_group_TF150" {
  description         = "test-terraform-TF150"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
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
