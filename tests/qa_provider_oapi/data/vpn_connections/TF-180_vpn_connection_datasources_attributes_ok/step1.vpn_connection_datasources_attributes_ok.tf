resource "outscale_virtual_gateway" "My_VGW" {
 connection_type = "ipsec.1"  
} 

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn          = 65000
    public_ip        = "198.18.7.207"
    connection_type  = "ipsec.1"
}

resource "outscale_client_gateway" "My_CGW_2" {
    bgp_asn          = 64900
    public_ip        = "198.18.7.206"
    connection_type  = "ipsec.1"
}

resource "outscale_vpn_connection" "vpn_connection_1" {
    client_gateway_id  = outscale_client_gateway.My_CGW.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type = "ipsec.1"
    static_routes_only = false
    tags {
        key   = "Project"
        value = "Terraform"
     }
    tags {
        key   = "Type"
        value = "Dynamic"
     }
}

resource "outscale_vpn_connection" "vpn_connection_2" {
    client_gateway_id  = outscale_client_gateway.My_CGW_2.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id    
    connection_type = "ipsec.1"
    static_routes_only = true
}

data "outscale_vpn_connections" "data_vpn_connections_1" {
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

data "outscale_vpn_connections" "data_vpn_connections_2" {
    filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection.vpn_connection_1.vpn_connection_id, outscale_vpn_connection.vpn_connection_2.vpn_connection_id]
    }
}
