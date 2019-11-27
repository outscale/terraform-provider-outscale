resource "outscale_vm" "outscale_vm_centos" {
    count = 2                                             # plus testWebsite one already created

    image_id               = "ami-be23e98b"
    vm_type                   = "c4.large"
    keypair_name           = "integ_sut_keypair"
    #firewall_rules_set_ids = ["sg-c73d3b6b"]
    security_group_ids     = ["sg-c73d3b6b"]    # tempo test
}

data "outscale_vms" "outscale_vms" {
   depends_on =["outscale_vm.outscale_vm_centos"]
   
   filter {
      name  = "vm_ids"
      values = [outscale_vm.outscale_vm_centos.0.vm_id,outscale_vm.outscale_vm_centos.1.vm_id]
   }  
}
