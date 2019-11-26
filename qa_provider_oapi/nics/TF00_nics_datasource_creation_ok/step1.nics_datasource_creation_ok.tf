resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    subregion_name = format("%s%s", var.region, "a")
    ip_range        = "10.0.0.0/16"
    net_id          = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}


data "outscale_nics" "outscale_nics" {
    network_interface_id = [outscale_nic.outscale_nic.nic_id, "outscale_nic.outscale_nic2.nic_id"]
}
