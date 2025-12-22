resource "outscale_virtual_gateway" "outscale_virtual_gateway2" {
 connection_type = "ipsec.1"
 tags {
  key = "name"
  value = "test-VGW-2"
 }
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_virtual_gateway_link" "outscale_virtual_gateway_link" {
    virtual_gateway_id = outscale_virtual_gateway.outscale_virtual_gateway2.virtual_gateway_id
    net_id              = outscale_net.outscale_net.net_id
}

