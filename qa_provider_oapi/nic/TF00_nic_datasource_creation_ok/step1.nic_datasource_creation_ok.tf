resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    availability_zone   = format("%s%s", var.region, "a")
    cidr_block          = "10.0.0.0/16"
    vpc_id              = outscale_lin.outscale_lin.vpc_id
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

data "outscale_nic" "outscale_nic" {
    network_interface_id = outscale_nic.outscale_nic.network_interface_id
    subnet_id            = outscale_subnet.oustcale_subnet.subnet_id
}
