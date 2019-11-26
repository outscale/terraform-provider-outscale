resource "outscale_lin" "outscale_lin" {
    count = 1

    ip_range          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    sub_region_name = format("%s%s", var.region, "a")
    ip_range        = "10.0.0.0/16"
    net_id          = outscale_lin.outscale_lin.net_id
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

resource "outscale_nic_link" "outscale_nic_link" {
    device_index            = 1
    #instance_id             = ""
    netword_interface_id    = outscale_nic.outscale_nic.network_interface_id
}
