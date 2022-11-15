resource "outscale_virtual_gateway" "My_VGW" { 
 connection_type = "ipsec.1"  
} 

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn          = 65000
    public_ip        = "198.18.7.207"
    connection_type  = "ipsec.1"
}


resource "outscale_vpn_connection" "outscale_vpn_connection" {
    client_gateway_id  = outscale_client_gateway.My_CGW.client_gateway_id
    virtual_gateway_id = outscale_virtual_gateway.My_VGW.virtual_gateway_id
    connection_type    = "ipsec.1"
    static_routes_only = true
    tags {
        key   = "Name"
        value = "test-VPN"
    }
}

resource "outscale_vpn_connection_route" "route2" {
 vpn_connection_id  = outscale_vpn_connection.outscale_vpn_connection.vpn_connection_id
 destination_ip_range = "40.0.0.0/16"
}
data "outscale_vpn_connection" "outscale_vpn_connection" {
filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection_route.route2.vpn_connection_id]
    }
filter {
       name   = "route_destination_ip_ranges"
       values = [outscale_vpn_connection_route.route2.destination_ip_range]
    }
}
