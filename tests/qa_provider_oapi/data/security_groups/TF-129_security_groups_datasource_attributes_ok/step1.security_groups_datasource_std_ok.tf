resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}


resource "outscale_security_group" "outscale_security_group" {
    description         = "test group-1"
    security_group_name = "terraform-TF125"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = ":outscale_net_resource2"
    }
    tags {
      key = "Key:"
      value = "value-tags"
     }
}


resource "outscale_security_group" "outscale_security_group2" {
    description         = "test group-2"
    security_group_name = "terraform-TF125-2"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_sg2"
    }
}

resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    from_port_range   = "80"
    to_port_range     = "80"
    ip_protocol       = "tcp"
    ip_range          = "46.231.147.88/32"
}

data "outscale_security_group" "outscale_security_groupd" {
    filter {   
        name   = "security_group_ids"
        values = [outscale_security_group.outscale_security_group.security_group_id]   
    }
depends_on=[outscale_security_group_rule.outscale_security_group_rule]
}

data "outscale_security_group" "outscale_security_group" {
    filter {
        name    = "account_ids"
        values  = [outscale_security_group.outscale_security_group.account_id]
    }
    filter {
        name    = "descriptions"
        values  = [outscale_security_group.outscale_security_group.description]
    }
depends_on=[outscale_security_group_rule.outscale_security_group_rule]
}

data "outscale_security_group" "filters-inbound" {
   filter {
        name    = "inbound_rule_account_ids"
        values  = [outscale_security_group.outscale_security_group.account_id]
   }
   filter {
        name    = "inbound_rule_from_port_ranges"
        values  = [80]
   }
   filter {
        name    = "inbound_rule_ip_ranges"
        values  = ["46.231.147.88/32"]
   }
   filter {
        name    = "inbound_rule_protocols"
        values  = ["tcp"]
   }
   filter {
        name    = "inbound_rule_to_port_ranges"
        values  = [80]
   }
depends_on=[outscale_security_group_rule.outscale_security_group_rule]
}



data "outscale_security_group" "filters-tags" {
   filter {
        name    = "tag_keys"
        values  = ["Key:"]
    }
   filter {
        name    = "tag_values"
        values  = [":outscale_net_resource2"]
    }
   filter {
        name    = "tags"
        values  = ["Name=:outscale_net_resource2"]
    }
depends_on=[outscale_security_group_rule.outscale_security_group_rule]
}

