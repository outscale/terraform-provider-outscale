resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    sub_region_name = format("%s%s", var.region, "a")
    ip_range        = "10.0.0.0/16"
    net_id          = outscale_net.outscale_net.net_id
}

output "outscale_subnet" {
    value = outscale_subnet.outscale_subnet.subnet_id
}
