resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF149"
}
resource "outscale_vm" "outscale_vm" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name
    placement_subregion_name = "${var.region}a"
}

data "outscale_vm_state" "outscale_vm_state" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id]
     }
}
