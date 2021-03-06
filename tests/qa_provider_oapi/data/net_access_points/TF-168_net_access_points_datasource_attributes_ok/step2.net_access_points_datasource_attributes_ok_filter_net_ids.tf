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
   route_table_ids = [outscale_route_table.route_table-1[1].route_table_id,outscale_route_table.route_table-2.route_table_id]
   service_name    = var.service_name
   tags {
          key      = "name"
          value    = "terraform-Net-Access-Point"
   }
   tags {
          key      = "test-terraform"
          value    = "Net-Access-Point"
   }
}

resource "outscale_net_access_point" "net_access_point_2" {
   net_id          = outscale_net.outscale_net.net_id
   route_table_ids = [outscale_route_table.route_table-1[0].route_table_id]
   service_name    = var.service_name
}

data "outscale_net_access_points" "data_net_access_point-2" {
    filter {
        name     = "net_ids"
        values   = [outscale_net_access_point.net_access_point_1.net_id]
      }
}

