resource "outscale_client_gateway" "outscale_client_gateway" {
    bgp_asn = random_integer.bgp_asn[0].result
    public_ip       = "171.33.75.123"
    connection_type = "ipsec.1"
    tags {
      key           = "name-terraform"
      value         = "CGW_1_terraform"
    }
    tags {
      key           = "project"
      value         = "terraform"
    }
}

resource "outscale_client_gateway" "outscale_client_gateway_2" {
    bgp_asn = random_integer.bgp_asn[1].result
    public_ip       = "171.33.75.43"
    connection_type = "ipsec.1"
    tags {
      key           = "name-terraform"
      value         = "CGW_2_terraform"
    }
}

data "outscale_client_gateways" "outscale_client_gateways" {
 filter {
       name   = "bgp_asns"
       values = [outscale_client_gateway.outscale_client_gateway.bgp_asn, outscale_client_gateway.outscale_client_gateway_2.bgp_asn]
    }
 filter {
       name   = "public_ips"
       values = [outscale_client_gateway.outscale_client_gateway.public_ip, outscale_client_gateway.outscale_client_gateway_2.public_ip]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway_2.state]

  }
depends_on=[outscale_client_gateway.outscale_client_gateway, outscale_client_gateway.outscale_client_gateway_2]
}

data "outscale_client_gateways" "outscale_client_gateways-2" {
 filter {
       name   = "tags"
       values = ["name-terraform=CGW_1_terraform"]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
  }
}

data "outscale_client_gateways" "outscale_client_gateways-3" {
  filter {
       name   = "tag_keys"
       values = ["name-terraform"]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
  }
}


data "outscale_client_gateways" "outscale_client_gateways-4" {
  filter {
       name   = "tag_values"
       values = ["CGW_1_terraform","CGW_2_terraform"]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
  }
}
