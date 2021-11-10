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
 tags {
        key   = "Name"
        value = ":outscale_net_resource2"
    }
    tags {
      key = "Key:"
      value = "value-tags"
     }
}

resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    from_port_range   = "80"
    to_port_range     = "80"
    ip_protocol       = "tcp"
    ip_range          = "46.231.147.8/32"
}


resource "outscale_security_group_rule" "outscale_security_group_rule-2" {
    flow              = "Outbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    rules {
     from_port_range   = "8080"
     to_port_range     = "8080"
     ip_protocol       = "tcp"
     ip_ranges         = ["46.231.147.8/32"]

     }
}


resource "outscale_security_group_rule" "outscale_security_group_rule-3" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    rules {
     from_port_range   = "22"
     to_port_range     = "22"
     ip_protocol       = "tcp"
      security_groups_members {
           account_id         = outscale_security_group.outscale_security_group2.account_id
           security_group_id  = outscale_security_group.outscale_security_group2.security_group_id
       }
     }
}

resource "outscale_security_group_rule" "outscale_security_group_rule-3_2" {
    flow              = "Outbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    rules {
     from_port_range   = "22"
     to_port_range     = "22"
     ip_protocol       = "tcp"
      security_groups_members {
           account_id         = outscale_security_group.outscale_security_group2.account_id
           security_group_id  = outscale_security_group.outscale_security_group2.security_group_id
       }
     }
}



data "outscale_security_groups" "outscale_security_groupsd" {
    filter {   
        name   = "security_group_ids"
        values = [outscale_security_group.outscale_security_group.security_group_id,outscale_security_group.outscale_security_group2.security_group_id]   
    }
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}

data "outscale_security_groups" "outscale_security_groups" {
    filter {
        name    = "account_ids"
        values  = [outscale_security_group.outscale_security_group.account_id]
    }
    filter {
        name    = "descriptions"
        values  = [outscale_security_group.outscale_security_group.description, outscale_security_group.outscale_security_group2.description]
    }
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}

data "outscale_security_groups" "filters-inbounds" {
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
        values  = ["46.231.147.8/32"]
   }
   filter {
        name    = "inbound_rule_protocols"
        values  = ["tcp"]
   }
   filter {
        name    = "inbound_rule_security_group_ids"
        values  = [outscale_security_group.outscale_security_group2.security_group_id]
   }
   filter {
        name    = "inbound_rule_security_group_names"
        values  = [outscale_security_group.outscale_security_group2.security_group_name]
   }
   filter {
        name    = "inbound_rule_to_port_ranges"
        values  = [80]
   }
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}

data "outscale_security_groups" "filters-outbounds" {
   filter {
        name    = "net_ids"
        values  = [outscale_net.outscale_net.net_id]
   }
filter {
        name    = "outbound_rule_account_ids"
        values  = [outscale_security_group.outscale_security_group.account_id]
   }
  filter {
        name    = "outbound_rule_from_port_ranges"
        values  =  [8080]
   }
  filter {
        name    = "outbound_rule_ip_ranges"
        values  = ["46.231.147.8/32"]
   }
  filter {
        name    = "outbound_rule_protocols"
        values  = ["tcp"]
   }
  filter {
        name    = "outbound_rule_security_group_ids"
        values  = [outscale_security_group.outscale_security_group2.security_group_id]
   }
   filter {
        name    = "outbound_rule_security_group_names"
        values  = [outscale_security_group.outscale_security_group2.security_group_name]
   }
   filter {
        name    = "outbound_rule_to_port_ranges"
        values  =  [8080]
   }
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}


data "outscale_security_groups" "filters-sgs" {
   filter {
        name    = "security_group_ids"
        values  = [outscale_security_group.outscale_security_group2.security_group_id]
   }
   filter {
        name    = "security_group_names"
        values  = [outscale_security_group.outscale_security_group2.security_group_name]
   }
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}


data "outscale_security_groups" "filters-tags" {
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
depends_on=[outscale_security_group_rule.outscale_security_group_rule,outscale_security_group_rule.outscale_security_group_rule-3,outscale_security_group_rule.outscale_security_group_rule-2,outscale_security_group_rule.outscale_security_group_rule-3_2]
}

