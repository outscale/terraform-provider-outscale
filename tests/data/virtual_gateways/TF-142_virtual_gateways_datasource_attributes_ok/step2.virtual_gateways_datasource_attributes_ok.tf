resource "outscale_virtual_gateway" "outscale_virtual_gateway" {
 connection_type = "ipsec.1"
 tags {
  key = "name-Datasources"
  value = "VGW-1-Datasources"
  }
 tags {
  key = "Project-Datasources"
  value = "Terraform-Datasources"
  }
}

resource "outscale_virtual_gateway" "outscale_virtual_gateway2" {
 connection_type = "ipsec.1"
 tags {
  key = "name"
  value = "VGW-2-Datasources"
  }
}


data "outscale_virtual_gateways" "outscale_virtual_gateways-3" {
filter {
        name   = "tags"
        values = ["name-Datasources=VGW-1-Datasources"]
    }
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id,outscale_virtual_gateway.outscale_virtual_gateway2.virtual_gateway_id]
    }
}
data "outscale_virtual_gateways" "outscale_virtual_gateways-4" {
filter {
        name   = "tag_keys"
        values = ["Project-Datasources"]
    }
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id,outscale_virtual_gateway.outscale_virtual_gateway2.virtual_gateway_id]
    }
}
data "outscale_virtual_gateways" "outscale_virtual_gateways-5" {
filter {
        name   = "tag_values"
        values = ["VGW-2-Datasources"]
    }
filter {
        name   = "virtual_gateway_ids"
        values = [outscale_virtual_gateway.outscale_virtual_gateway.virtual_gateway_id,outscale_virtual_gateway.outscale_virtual_gateway2.virtual_gateway_id]
    }
}

