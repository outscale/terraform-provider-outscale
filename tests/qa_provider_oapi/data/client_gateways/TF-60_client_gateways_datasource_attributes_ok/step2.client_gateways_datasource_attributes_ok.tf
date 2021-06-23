resource "outscale_client_gateway" "outscale_client_gateway" {
    bgp_asn     = 571
    public_ip  = "171.33.75.123"
    connection_type        = "ipsec.1"
    tags {
     key = "name-mzi"
     value = "CGW_1_mzi"
    }

 tags {
     key = "project"
     value = "terraform"
    }
}

resource "outscale_client_gateway" "outscale_client_gateway_2" {
    bgp_asn     = 575
    public_ip  = "171.33.75.43"
    connection_type        = "ipsec.1"
    tags {
     key = "name-mzi"
     value = "CGW_2_mzi"
    }
}

data "outscale_client_gateways" "outscale_client_gateways" {
filter {
       name   = "client_gateway_ids"
       values = [outscale_client_gateway.outscale_client_gateway.client_gateway_id,outscale_client_gateway.outscale_client_gateway_2.client_gateway_id]
    }
filter {
       name   = "bgp_asns"
       values = [outscale_client_gateway.outscale_client_gateway.bgp_asn]
    }
filter {
       name   = "public_ips"
       values = [outscale_client_gateway.outscale_client_gateway.public_ip]
    }
filter {
       name   = "tags"
       values = ["name-mzi=CGW_1_mzi"]
    }
filter {
       name   = "tag_keys"
       values = ["name-mzi"]
    }
filter {
       name   = "tag_values"
       values = ["CGW_1_mzi"]
    }
}

