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



resource "outscale_security_group_rule" "outscale_security_group_rule-2" {
    flow              = "Outbound"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
    rules {
     from_port_range   = "8080"
     to_port_range     = "8080"
     ip_protocol       = "tcp"
     ip_ranges         = ["46.231.147.88/32"]

     }
}


data "outscale_security_group" "filters-outbound" {
   filter {
        name    = "net_ids"
        values  = [outscale_net.outscale_net.net_id]
   }
  filter {
        name    = "outbound_rule_from_port_ranges"
        values  =  [8080]
   }
  filter {
        name    = "outbound_rule_ip_ranges"
        values  = ["46.231.147.88/32"]
   }
  filter {
        name    = "outbound_rule_protocols"
        values  = ["tcp"]
   }
   filter {
        name    = "outbound_rule_to_port_ranges"
        values  =  [8080]
   }
depends_on=[outscale_security_group_rule.outscale_security_group_rule-2]
}
