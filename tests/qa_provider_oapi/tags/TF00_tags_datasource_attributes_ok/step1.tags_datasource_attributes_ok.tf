resource "outscale_vm" "outscale_vm" {
    image_id               = var.image_id
    vm_type                = var.vm_type
    keypair_name           = var.keypair_name
    security_group_ids     = [var.security_group_id]
    tags  {                              
        key   = "name7"
        value = "testDataSource7"        
    }                           
}

resource "outscale_vm" "outscale_vm2" {
    image_id               = var.image_id
    vm_type                = var.vm_type
    keypair_name           = var.keypair_name
    security_group_ids     = [var.security_group_id]
    tags  {                              
        key   = "name7"
        value = "testDataSource72"        
    }                 
}

data "outscale_tags" "outscale_tags" {
    filter {
        name = "resource_ids"
        values = [outscale_vm.outscale_vm.vm_id, outscale_vm.outscale_vm2.vm_id]
    }
}
