resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
connection_type = "ipsec.1"
 tags {
  key = "Project-Datasource"
  value = "Terraform-Datasource"
  }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_virtual_gateway_link" "outscale_virtual_gateway_link" {
    virtual_gateway_id = outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id
    net_id              = outscale_net.outscale_net.net_id
}

data "outscale_virtual_gateway" "outscale_vpn_gateway" {
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id]
    }
}
