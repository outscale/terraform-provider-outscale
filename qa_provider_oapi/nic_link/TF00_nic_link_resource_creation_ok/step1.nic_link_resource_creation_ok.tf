resource "outscale_vm" "outscale_instance" {                  # OK
    count = 1

    image_id               = "ami-880caa66"
    #image_id              = "ami-5ad76458"
    instance_type          = "c4.large"
    #instance_type         = "c4.large"
    keypair_name           = "integ_sut_keypair"
    firewall_rules_set_ids = ["sg-c73d3b6b"]
    

    provisioner "local-exec" {
        command = "date"
    }
    provisioner "local-exec" {
        command = "date; echo ${self.image_id} ${self.instance_type} ${self.id
    }
}

resource "outscale_net" "outscale_net" {
    count = 1

    ip_range          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    sub_region_name = format("%s%s", var.region, "a")
    ip_range        = "10.0.0.0/16"
    net_id          = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_nic_link" "outscale_nic_link" {
    device_index            = "1"
    instance_id             = outscale_vm.outscale_vm.instance_id
    netword_interface_id    = outscale_nic.outscale_nic.network_interface_id
}
