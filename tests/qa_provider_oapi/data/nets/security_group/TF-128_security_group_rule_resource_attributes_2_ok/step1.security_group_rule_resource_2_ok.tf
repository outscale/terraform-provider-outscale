resource "outscale_net" "outscale_net" {
   ip_range = "10.0.0.0/16"
    tags {
        key   = "Name"
        value = "outscale_net_resource2"
    }
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-terraform-test"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_sg"
    }
}

resource "outscale_security_group" "outscale_security_group2" {
    description         = "test group"
    security_group_name = "sg2-terraform-test"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_sg2"
    }
}

resource "outscale_security_group_rule" "outscale_security_group_rule-3" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.outscale_security_group.id
    rules {
     from_port_range   = "22"
     to_port_range     = "22"
     ip_protocol       = "tcp"
      security_groups_members {
           account_id         = outscale_security_group.outscale_security_group2.account_id
           security_group_id  = outscale_security_group.outscale_security_group2.id
       }
     }
depends_on = [outscale_security_group.outscale_security_group2]
}

resource "outscale_security_group_rule" "outscale_security_group_rule-3_2" {
    flow              = "Outbound"
    security_group_id = outscale_security_group.outscale_security_group.id
    rules {
     from_port_range   = "22"
     to_port_range     = "22"
     ip_protocol       = "tcp"
      security_groups_members {
           account_id         = outscale_security_group.outscale_security_group2.account_id
           security_group_name  = outscale_security_group.outscale_security_group2.security_group_name
       }
     }
depends_on = [outscale_security_group.outscale_security_group2, outscale_security_group_rule.outscale_security_group_rule-3]
}

