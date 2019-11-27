resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
}

resource "outscale_route_table" "outscale_route_table2" {
    net_id = outscale_net.outscale_net.net_id
}

data "outscale_route_tables" "outscale_route_tables" {
    filter {
        name   = "route_table_ids"
        values = [outscale_route_table.outscale_route_table.route_table_id, outscale_route_table.outscale_route_table2.route_table_id]
    }
}
