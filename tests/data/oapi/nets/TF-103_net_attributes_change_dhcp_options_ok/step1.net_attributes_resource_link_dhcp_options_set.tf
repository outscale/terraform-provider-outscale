resource "outscale_dhcp_option" "outscale_dhcp_option1" {
  domain_name = "test123.fr"
  tags {
    key   = "name-1"
    value = "test-MZI-1"
  }
}

resource "outscale_net" "outscale_net" {
  ip_range = "10.0.0.0/16"
  tags {
    key   = "name"
    value = "test-net-attributes"
  }
}
resource "outscale_net_attributes" "outscale_net_attributes" {
  net_id              = outscale_net.outscale_net.net_id
  dhcp_options_set_id = outscale_dhcp_option.outscale_dhcp_option1.id
}
