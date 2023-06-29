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

data "outscale_virtual_gateway" "outscale_virtual_gateways-2" {
filter {
        name   = "link_net_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.net_id]
    }
}

data "outscale_virtual_gateway" "outscale_vpn_gateway" {
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id]
    }
}

data "outscale_virtual_gateway" "outscale_virtual_gateways-3" {
 filter {
        name   = "tags"
        values = ["Project-Datasource=Terraform-Datasource"]
    }
 filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id]
    }
depends_on = [outscale_virtual_gateway_link.outscale_virtual_gateway_link]
}

data "outscale_virtual_gateway" "outscale_virtual_gateways-4" {
 filter {
        name   = "tag_keys"
        values = ["Project-Datasource"]
    }
 filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id]
    }
depends_on = [outscale_virtual_gateway_link.outscale_virtual_gateway_link]
}

data "outscale_virtual_gateway" "outscale_virtual_gateways-5" {
filter {
        name   = "tag_values"
        values = ["Terraform-Datasource"]
    }
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway_link.outscale_virtual_gateway_link.virtual_gateway_id]
    }
depends_on = [outscale_virtual_gateway_link.outscale_virtual_gateway_link]
}

