resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF111"
}
resource "outscale_vm" "outscale_vm" {
    image_id                 = var.image_id
    vm_type                  = var.vm_type
    keypair_name             = outscale_keypair.my_keypair.keypair_name
    security_group_ids       = [outscale_security_group.outscale_security_group.id]
    subnet_id                = outscale_subnet.outscale_subnet.id
}


resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name = "${var.region}a"
    ip_range       = "10.0.0.0/16"
    net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_nic_link" "outscale_nic_link" {
    device_number = "1"
    vm_id         = outscale_vm.outscale_vm.vm_id
    nic_id        = outscale_nic.outscale_nic.nic_id
}
