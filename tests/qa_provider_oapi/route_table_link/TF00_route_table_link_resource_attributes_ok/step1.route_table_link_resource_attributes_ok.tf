resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "terraform-RT"
    }
}

resource "outscale_subnet" "outscale_subnet" {
    net_id   = outscale_net.outscale_net.net_id
    ip_range = "10.0.0.0/18"
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    route_table_id  = outscale_route_table.outscale_route_table.route_table_id
    subnet_id       = outscale_subnet.outscale_subnet.subnet_id
}
