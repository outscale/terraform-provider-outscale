resource "outscale_virtual_gateway" "My_VGW" {
 connection_type = "ipsec.1"
}

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn = random_integer.bgp_asn[0].result
    public_ip        = "198.18.7.207"
    connection_type  = "ipsec.1"
}

resource "outscale_client_gateway" "My_CGW_2" {
    bgp_asn = random_integer.bgp_asn[1].result
    public_ip        = "198.18.7.205"
    connection_type  = "ipsec.1"
}


resource "outscale_vpn_connection" "dynamic_vpn_connection" {
    client_gateway_id  = outscale_client_gateway.My_CGW_2.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type = "ipsec.1"
    static_routes_only = false
    tags {
        key   = "Type"
        value = "Dynamic"
    }
    tags {
        key   = "Project"
        value = "Terraform"
    }
}
