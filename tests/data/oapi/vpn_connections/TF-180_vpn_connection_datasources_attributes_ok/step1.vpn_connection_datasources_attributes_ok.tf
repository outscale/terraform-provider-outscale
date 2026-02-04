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
    public_ip        = "198.18.7.206"
    connection_type  = "ipsec.1"
}

resource "outscale_vpn_connection" "vpn_connection_1" {
    client_gateway_id  = outscale_client_gateway.My_CGW.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type = "ipsec.1"
    static_routes_only = false
    tags {
        key   = "Test-TF180"
        value = "Terraform"
     }
    tags {
        key   = "Type:TF-180"
        value = "Dynamic-TF:180"
     }
}

resource "outscale_vpn_connection" "vpn_connection_2" {
    client_gateway_id  = outscale_client_gateway.My_CGW_2.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type = "ipsec.1"
    static_routes_only = true
    tags {
        key   = "Type:TF-180"
        value = "Static-TF:180"
     }
}
