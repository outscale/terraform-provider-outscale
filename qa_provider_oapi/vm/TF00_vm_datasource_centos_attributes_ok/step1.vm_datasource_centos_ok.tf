resource "outscale_vm" "outscale_vm_centos" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = var.keypair_name
    security_group_names = [var.security_group_name]
}

data "outscale_vm" "outscale_vm" {
	filter {
        name   = "vm_ids"
        #values = [outscale_vm.outscale_vm_centos.id]
        values = [outscale_vm.outscale_vm_centos.vm_id]
        #vm_id = outscale_vm.outscale_vm_centos.vm.0.vm_id
        #vm_id = outscale_vm.outscale_vm_centos.vm.0.id
    }
}

resource "outscale_vm" "outscale_vm_centos2" {
    count = 2

    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = var.keypair_name
    security_group_ids = [var.security_group_id]
}

data "outscale_vm" "outscale_vm_centos2_0" {
	filter {
        name   = "vm_ids"
        #values = [outscale_vm.outscale_vm_centos2.0.vm_id]
        values = [outscale_vm.outscale_vm_centos2.*.vm_id[0]]
    }
}

data "outscale_vm" "outscale_vm_centos2_1" {
	filter {
        name   = "vm_ids"
        #values = [outscale_vm.outscale_vm_centos2.1.vm_id]
        values = [outscale_vm.outscale_vm_centos2.*.vm_id[1]]
    }
}

# TODO
# Not yet adapted to oAPI

#output "datasource_arch" {
#  value = data.outscale_vm.outscale_vm.instances_set.0.architecture
#}

#output "datasource_network" {
#  value = data.outscale_vm.outscale_vm.instances_set.0.network_interface_set.0.network_interface_id
#}
