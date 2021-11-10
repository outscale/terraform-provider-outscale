resource "outscale_client_gateway" "outscale_client_gateway" {
    bgp_asn     = 571
    public_ip  = "171.33.75.123"
    connection_type        = "ipsec.1"
    tags {
     key = "name:mzi"
     value = "CGW_1:mzi"
    }
 tags {
     key = "project"
     value = "terraform"
    }
}


data "outscale_client_gateway" "outscale_client_gateway_2" {
 filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
    }

 filter {
       name   = "bgp_asns"
       values = [outscale_client_gateway.outscale_client_gateway.bgp_asn]
    }
 filter {
       name   = "public_ips"
       values = [outscale_client_gateway.outscale_client_gateway.public_ip]
    }
}

data "outscale_client_gateway" "outscale_client_gateway_3" {
 filter {
       name   = "tags"
       values = ["name:mzi=CGW_1:mzi"]
    }
 filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
    }
depends_on=[outscale_client_gateway.outscale_client_gateway]
}

data "outscale_client_gateway" "outscale_client_gateway_4" {
 filter {
       name   = "tag_keys"
       values = ["name:mzi"]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
    }
depends_on=[outscale_client_gateway.outscale_client_gateway]
}

data "outscale_client_gateway" "outscale_client_gateway_5" {
 filter {
       name   = "tag_values"
       values = ["CGW_1:mzi"]
    }
  filter {
       name   = "states"
       values = [outscale_client_gateway.outscale_client_gateway.state]
    }
depends_on=[outscale_client_gateway.outscale_client_gateway]
}
