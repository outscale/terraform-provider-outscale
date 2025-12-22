resource "outscale_virtual_gateway" "My_VGW" {
 connection_type = "ipsec.1"
}

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn = random_integer.bgp_asn[0].result
    public_ip        = "198.18.7.207"
    connection_type  = "ipsec.1"
}

resource "outscale_vpn_connection" "static_vpn_connection" {
    client_gateway_id  = outscale_client_gateway.My_CGW.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type = "ipsec.1"
    static_routes_only = true
    tags {
        key   = "Type"
        value = "Static"
    }
    tags {
        key   = "Project"
        value = "Terraform"
    }
}
