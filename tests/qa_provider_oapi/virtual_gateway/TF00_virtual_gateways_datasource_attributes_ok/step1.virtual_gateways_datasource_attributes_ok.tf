resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
 connection_type = "ipsec.1"
 tags {
  key = "name"
  value = "test-VGW-1"
  }
}

resource "outscale_virtual_gateway" "outscale_virtual_gateway2" {
 connection_type = "ipsec.1"
 tags {
  key = "name"
  value = "test-VGW-2"
  }
}

data "outscale_virtual_gateways" "outscale_virtual_gateways" {
    virtual_gateway_id = [ outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id,  outscale_virtual_gateway.outscale_virtual_gateway2.virtual_gateway_id]
}