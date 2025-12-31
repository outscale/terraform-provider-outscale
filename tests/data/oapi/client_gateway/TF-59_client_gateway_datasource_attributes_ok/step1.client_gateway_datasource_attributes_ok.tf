resource "outscale_client_gateway" "outscale_client_gateway" {
    bgp_asn = random_integer.bgp_asn[0].result
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



data "outscale_client_gateway" "outscale_client_gateway_2" {
filter {
       name   = "client_gateway_ids"
       values = [outscale_client_gateway.outscale_client_gateway.client_gateway_id]
    }
}
