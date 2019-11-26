resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    sub_region_name = format("%s%s", var.region, "a")
    ip_range        = "10.0.0.0/17"
    net_id          = outscale_lin.outscale_lin.vpc_id
}
