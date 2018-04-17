resource "outscale_lin" "outscale_lin" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_subnet" "outscale_subnet" {
  vpc_id     = "${outscale_lin.outscale_lin.vpc_id}"
  cidr_block = "10.0.0.0/18"
}

resource "outscale_route_table_link" "outscale_route_table_link" {
  route_table_id = "${outscale_route_table.outscale_route_table.route_table_id}"
  subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_client_endpoint" "outscale_client_endpoint" {
  bgp_asn    = "3"
  ip_address = "171.33.74.122"
  type       = "ipsec.1"
}

resource "outscale_dhcp_option" "outscale_dhcp_option" {}

resource "outscale_dhcp_option_link" "outscale_dhcp_option_link" {
  dhcp_options_id = "${outscale_dhcp_option.outscale_dhcp_option.dhcp_options_id}"
  vpc_id          = "${outscale_lin.outscale_lin.vpc_id}"
}
