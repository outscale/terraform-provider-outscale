resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test-r"
tags {
      key = "Key"
      value = "value-tags"
     }
tags {
      key = "Key-2"
      value = "value-tags-2"
     }
}

resource "outscale_security_group" "outscale_security_group2" {
    description         = "test group"
    security_group_name = "sg1-test-group_test-r2"
}

data "outscale_security_groups" "outscale_security_groups" {
    filter {
        name  = "security_group_ids"
        values = [outscale_security_group.outscale_security_group.security_group_id,outscale_security_group.outscale_security_group2.security_group_id]
    }
}
