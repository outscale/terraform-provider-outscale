resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_security_group" "security_group_TF151" {
  description         = "test-terraform-TF151"
  security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "outscale_vm_centos" {
    count = 2                                             # plus testWebsite one already created

     image_id               = var.image_id
     vm_type                = var.vm_type
     keypair_name           = outscale_keypair.my_keypair.keypair_name
     security_group_ids       = [outscale_security_group.security_group_TF151.security_group_id]
}

data "outscale_vms" "outscale_vms" {
   depends_on =[outscale_vm.outscale_vm_centos]

   filter {
      name  = "vm_ids"
      values = [outscale_vm.outscale_vm_centos.0.vm_id,outscale_vm.outscale_vm_centos.1.vm_id]
   }
}

resource "outscale_vm" "outscale_vm_tags" {
     image_id               = var.image_id
     vm_type                = var.vm_type
     keypair_name           = outscale_keypair.my_keypair.keypair_name
     security_group_ids       = [outscale_security_group.security_group_TF151.security_group_id]
     tags {
        key = "name-B"
        value = "test-B"
     }
     tags {
         key = "Key"
         value = "value-tags"
     }
}
