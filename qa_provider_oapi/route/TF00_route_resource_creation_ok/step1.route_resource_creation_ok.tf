resource "outscale_net" "outscale_net" {
    count = 1

    cidr_block = "10.0.0.0/24"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    net_id = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_net_internet_gateway" "outscale_net_internet_gateway" {
    count = 1
}

resource "outscale_net_internet_gateway_link" "outscale_net_internet_gateway_link" {

    internet_gateway_id = "${outscale_net_internet_gateway.outscale_net_internet_gateway.internet_gateway_id}"
    net_id              = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_route" "outscale_route" {
    count = 1

    gateway_id           = "${outscale_net_internet_gateway.outscale_net_internet_gateway.internet_gateway_id}"
    destination_ip_range = "10.0.0.0/16"
    route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
}