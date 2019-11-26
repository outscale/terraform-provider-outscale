/*
resource "outscale_net" "outscale_net" {
    count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    net_id = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_route" "outscale_route" {
    count = 1

    destination_ip_range = "${outscale_net.outscale_net.ip_range}"
    route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
}
*/