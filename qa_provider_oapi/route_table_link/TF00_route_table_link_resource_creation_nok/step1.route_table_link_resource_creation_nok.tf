resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_net" "outscale_net1" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    net_id = "${outscale_net.outscale_net1.net_id}"

}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    net_id = "${outscale_net.outscale_net.net_id}"
    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    count = 1

    route_table_id  = "${outscale_route_table.outscale_route_table.route_table_id}"
    subnet_id       = "${outscale_subnet.outscale_subnet.subnet_id}"
}