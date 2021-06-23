resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF150"
}
resource "outscale_vm" "outscale_vm" {
   image_id                 = var.image_id
   vm_type                  = var.vm_type
   keypair_name             = outscale_keypair.my_keypair.keypair_name
}

resource "outscale_vm" "outscale_vm2" {
    image_id                  = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name 
}

data "outscale_vm_states" "outscale_vm_states" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id, outscale_vm.outscale_vm2.vm_id]
     }
}
