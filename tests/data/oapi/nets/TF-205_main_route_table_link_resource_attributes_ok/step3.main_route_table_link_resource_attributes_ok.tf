resource "outscale_net" "outscale_net" {
  ip_range = "10.0.0.0/16"
  tags {
    key = "name"
    value = "terraform-TF-117"
  }
}

resource "outscale_subnet" "outscale_subnet_1" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.0.0.0/24"
  tags {
    key = "name"
    value = "terraform-TF-117"
  }
}

resource "outscale_subnet" "outscale_subnet_2" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.0.1.0/24"
  tags {
    key = "name"
    value = "terraform-TF-117"
  }
}

resource "outscale_route_table" "outscale_route_table" {
  net_id = outscale_net.outscale_net.net_id
  tags {
    key = "name"
    value = "terraform-TF-117"
  }
}

resource "outscale_route_table_link" "outscale_route_table_link" {
  route_table_id  = outscale_route_table.outscale_route_table.route_table_id
  subnet_id       = outscale_subnet.outscale_subnet_1.subnet_id
}

resource "outscale_route_table_link" "outscale_route_table_link_2" {
  route_table_id  = outscale_route_table.outscale_route_table.route_table_id
  subnet_id       = outscale_subnet.outscale_subnet_2.subnet_id
}

resource "outscale_main_route_table_link" "main" {
  net_id   = outscale_net.outscale_net.net_id
  route_table_id = outscale_route_table.outscale_route_table.route_table_id
}
