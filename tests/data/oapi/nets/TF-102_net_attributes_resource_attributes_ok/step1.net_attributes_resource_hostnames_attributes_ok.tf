resource "outscale_dhcp_option" "dhcp_option_1" {
domain_name ="test"
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_net_attributes" "outscale_net_attributes" {
     net_id              = outscale_net.outscale_net.net_id
     dhcp_options_set_id = outscale_dhcp_option.dhcp_option_1.dhcp_options_set_id
}
