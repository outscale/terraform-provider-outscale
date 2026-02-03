resource "outscale_net" "outscale_net" {
     ip_range      = "10.0.0.0/16"
}

resource "outscale_route_table" "route_table-1" {
   count           = 2
   net_id          = outscale_net.outscale_net.net_id
}

resource "outscale_route_table" "route_table-2" {
   net_id          = outscale_net.outscale_net.net_id
}


resource "outscale_net_access_point" "net_access_point_1" {
   net_id          = outscale_net.outscale_net.net_id
   route_table_ids = [outscale_route_table.route_table-1[0].route_table_id, outscale_route_table.route_table-1[1].route_table_id, outscale_route_table.route_table-2.route_table_id]
   service_name    = "com.outscale.${var.region}.api"
}
