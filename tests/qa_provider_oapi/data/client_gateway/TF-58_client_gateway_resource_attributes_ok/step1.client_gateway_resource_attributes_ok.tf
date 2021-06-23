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
