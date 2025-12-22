resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF143"
}
resource "outscale_security_group" "public_sg" {
    description         = "test vms"
    security_group_name = "terraform-public-sg"
}

resource "outscale_vm" "outscale_vm_centos" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_names = [outscale_security_group.public_sg.security_group_name]
    tags {
     key = "name-1"
     value ="test-VM-tag"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_vm" "outscale_vm" {
  filter {
        name   = "vm_ids"
        values = [outscale_vm.outscale_vm_centos.vm_id]
 }
}

data "outscale_vm" "outscale_vm_2" {
filter {
  name   = "tags"
  values = ["name-1=test-VM-tag"]
 }
}

data "outscale_vm" "outscale_vm_3" {

filter {
  name   = "tag_keys"
  values = ["name-1"]
 }

filter {
  name   = "tag_values"
  values = ["test-VM-tag"]
 }
}

resource "outscale_vm" "outscale_vm_centos2" {
    count = 2

    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
    security_group_ids = [outscale_security_group.public_sg.security_group_id]
     tags {
     key = "name"
     value ="test-VM-tag-2"
    }
}

data "outscale_vm" "outscale_vm_centos2_0" {
	filter {
        name   = "vm_ids"
        values = [outscale_vm.outscale_vm_centos2.*.vm_id[0]]
    }
}

data "outscale_vm" "outscale_vm_centos2_1" {
	filter {
        name   = "vm_ids"
        values = [outscale_vm.outscale_vm_centos2.*.vm_id[1]]
    }
}
