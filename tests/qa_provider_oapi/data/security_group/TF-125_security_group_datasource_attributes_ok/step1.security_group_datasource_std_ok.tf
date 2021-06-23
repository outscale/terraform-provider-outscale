resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"

    tags {
        key   = "Name"
        value = "outscale_net_resource2"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_security_group" "outscale_security_groupd" {     
    filter {
        name   = "security_group_ids"
        values = [outscale_security_group.outscale_security_group.security_group_id]
        #values = [outscale_security_group.outscale_security_group.id]                 # test purposes only
    }
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test-d"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_security_group_net"
    }
}
