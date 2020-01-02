resource "outscale_net" "outscale_net" {
   # count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
   # count = 1

    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.10.0/24"
    net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "outscale_subnet2" {
    # count = 1

    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.20.0/24"
    net_id         = outscale_net.outscale_net.net_id
}
data "outscale_subnets" "outscale_subnets" {
    filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id, outscale_subnet.outscale_subnet2.subnet_id]
    }
}
