resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
   tags {
     key = "name"
     value = "terraform-RT"
    }
   tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_route_table" "outscale_route_table" {
    filter {
        name   = "route_table_ids"
        values = [outscale_route_table.outscale_route_table.route_table_id]
    }
}
