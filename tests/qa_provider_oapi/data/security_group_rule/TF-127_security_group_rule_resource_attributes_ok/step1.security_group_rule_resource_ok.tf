resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test"
}

resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.outscale_security_group.id
    from_port_range   = "0"
    to_port_range     = "0"
    ip_protocol       = "tcp"
    ip_range          = "0.0.0.0/0"
}

resource "outscale_security_group_rule" "outscale_security_group_rule_2" {
    flow              = "Inbound"
    from_port_range   = 22
    to_port_range     = 22
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
depends_on = [outscale_security_group_rule.outscale_security_group_rule]
}
