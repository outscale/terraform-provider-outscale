resource "outscale_route_table" "eth0" {
  net_id = outscale_net.my_net.net_id
}

resource "outscale_route_table_link" "my_route_table_link" {
  subnet_id      = outscale_subnet.public.subnet_id
  route_table_id = outscale_route_table.eth0.route_table_id
}

resource "outscale_route" "my_default_route" {
  destination_ip_range = "0.0.0.0/0"
  gateway_id           = outscale_internet_service.my_internet_service.internet_service_id
  route_table_id       = outscale_route_table.eth0.route_table_id
}