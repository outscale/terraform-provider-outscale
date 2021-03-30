resource "outscale_vm" "outscale_vm_centos" {
    count = 2                                             # plus testWebsite one already created

    image_id               = var.image_id
     vm_type                = var.vm_type
     keypair_name           = var.keypair_name
     security_group_ids     = [var.security_group_id]
}

data "outscale_vms" "outscale_vms" {
   depends_on =["outscale_vm.outscale_vm_centos"]
   
   filter {
      name  = "vm_ids"
      values = [outscale_vm.outscale_vm_centos.0.vm_id,outscale_vm.outscale_vm_centos.1.vm_id]
   }  
}
