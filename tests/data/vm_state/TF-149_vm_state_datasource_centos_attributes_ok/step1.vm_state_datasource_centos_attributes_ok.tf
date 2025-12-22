resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF149"
}

resource "outscale_security_group" "security_group_TF149" {
  description         = "test-terraform-TF149"
  security_group_name = "terraform-sg-149"
}

resource "outscale_vm" "outscale_vm" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name
    security_group_ids       = [outscale_security_group.security_group_TF149.security_group_id]
    placement_subregion_name = "${var.region}a"
}

data "outscale_vm_state" "outscale_vm_state" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id]
     }
}
