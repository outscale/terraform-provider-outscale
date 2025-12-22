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


data "outscale_vpn_connections" "data_vpn_connections_1" {
    filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection.vpn_connection_1.vpn_connection_id, outscale_vpn_connection.vpn_connection_2.vpn_connection_id]
    }
}

data "outscale_vpn_connections" "data_vpn_connections_2" {
    filter {
       name   = "client_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.client_gateway_id]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connections" "data_vpn_connections_3" {
    filter {
       name   = "virtual_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.virtual_gateway_id]
    }

   filter {
       name   = "static_routes_only"
       values = ["false"]
    }

depends_on =[outscale_vpn_connection.vpn_connection_1]
}


data "outscale_vpn_connections" "data_vpn_connections_4" {
filter {
       name   = "tag_keys"
       values = ["Test-TF180"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connections" "data_vpn_connections_5" {
filter {
       name   = "tag_values"
       values = ["Terraform"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connections" "data_vpn_connections_6" {
filter {
       name   = "tags"
       values = ["Type:TF-180=Dynamic-TF:180"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

