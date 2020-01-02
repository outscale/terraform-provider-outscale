resource "outscale_vm" "outscale_vm" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = var.keypair_name
    security_group_ids       = [var.security_group_id]
    #placement_subregion_name = format("%s%s", var.region, "a")
    #placement_tenancy        = "default"
}

data "outscale_vm_state" "outscale_vm_state" {
     filter {
         name   = "vm_ids"
         values = [outscale_vm.outscale_vm.vm_id]
     }
}
