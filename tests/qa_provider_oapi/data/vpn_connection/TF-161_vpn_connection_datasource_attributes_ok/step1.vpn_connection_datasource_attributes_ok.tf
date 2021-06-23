resource "outscale_virtual_gateway" "My_VGW" {
 connection_type = "ipsec.1"  
} 

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn          = 65000
    public_ip        = "198.18.7.207"
    connection_type  = "ipsec.1"
}

resource "outscale_vpn_connection" "vpn_connection_1" {
    client_gateway_id  = outscale_client_gateway.My_CGW.client_gateway_id
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

data "outscale_vpn_connection" "data_vpn_connection_1" {
    filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection.vpn_connection_1.vpn_connection_id]
    }

    filter {
       name   = "client_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.client_gateway_id]
    }

    filter {
       name   = "virtual_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.virtual_gateway_id]
    }

   filter {
       name   = "static_routes_only"
       values = [outscale_vpn_connection.vpn_connection_1.static_routes_only]
    }
}
