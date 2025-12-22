resource "outscale_virtual_gateway" "My_VGW" {
 connection_type = "ipsec.1"
}

resource "outscale_client_gateway" "My_CGW" {
    bgp_asn = random_integer.bgp_asn[0].result
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
        value = "Dynamic-TF161"
    }
    tags {
        key   = "Test"
        value = "Terraform-TF161"
    }
}

data "outscale_vpn_connection" "data_vpn_connection_1" {
    filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection.vpn_connection_1.vpn_connection_id]
    }
}
data "outscale_vpn_connection" "data_vpn_connection_2" {
    filter {
       name   = "client_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.client_gateway_id]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connection" "data_vpn_connection_3" {
    filter {
       name   = "virtual_gateway_ids"
       values = [outscale_vpn_connection.vpn_connection_1.virtual_gateway_id]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connection" "data_vpn_connection_4" {

   filter {
       name   = "static_routes_only"
       values = [outscale_vpn_connection.vpn_connection_1.static_routes_only]
    }
   filter {
       name   = "states"
       values = ["available"]
    }
   filter {
       name   = "vpn_connection_ids"
       values = [outscale_vpn_connection.vpn_connection_1.vpn_connection_id]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connection" "data_vpn_connection_5" {
filter {
       name   = "tag_keys"
       values = ["Type"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connection" "data_vpn_connection_6" {
filter {
       name   = "tag_values"
       values = ["Dynamic-TF161"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}

data "outscale_vpn_connection" "data_vpn_connection_7" {
filter {
       name   = "tags"
       values = ["Type=Dynamic-TF161"]
    }
filter {
       name   = "states"
       values = ["available"]
    }
depends_on =[outscale_vpn_connection.vpn_connection_1]
}
