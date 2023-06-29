resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name = "${var.region}a"
    ip_range       = "10.0.0.0/16"
    net_id         = outscale_net.outscale_net.net_id
    tags {
    key = "name"
    value = "terraform-subnet"
    }
}

output "outscale_subnet" {
    value = outscale_subnet.outscale_subnet.subnet_id
}
