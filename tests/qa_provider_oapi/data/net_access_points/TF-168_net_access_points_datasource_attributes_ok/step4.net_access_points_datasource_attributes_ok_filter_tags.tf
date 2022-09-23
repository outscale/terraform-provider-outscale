resource "outscale_net" "outscale_net" {
     ip_range      = "10.0.0.0/16"
}


resource "outscale_net" "outscale_net_2" {
     ip_range      = "10.10.0.0/16"
}


resource "outscale_route_table" "route_table-1" {
   net_id          = outscale_net.outscale_net.net_id
}

resource "outscale_route_table" "route_table-2" {
   net_id          = outscale_net.outscale_net_2.net_id
}

resource "outscale_net_access_point" "net_access_point_1" {
   net_id          = outscale_net.outscale_net.net_id
   route_table_ids = [outscale_route_table.route_table-1.route_table_id]
   service_name    = var.service_name
   tags {
          key      = "name"
          value    = "terraform-Net-Access-Point"
   }
}

resource "outscale_net_access_point" "net_access_point_2" {
   net_id          = outscale_net.outscale_net_2.net_id
   route_table_ids = [outscale_route_table.route_table-2.route_table_id]
   service_name    = var.service_name
tags {
          key      = "test-terraform"
          value    = "Net-Access-Point"
   }
}

data "outscale_net_access_points" "data_net_access_points_filter_tags" {

    filter {
        name     = "tags"
        values   = ["test-terraform=Net-Access-Point"]
      }

depends_on =[outscale_net_access_point.net_access_point_1]

}
