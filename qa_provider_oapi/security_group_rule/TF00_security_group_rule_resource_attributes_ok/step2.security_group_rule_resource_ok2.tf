resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    from_port_range   = 22
    to_port_range     = 22
    ip_protocol       = "tcp"
    ip_range          = "46.231.147.2/32"
    security_group_id = "${outscale_security_group.outscale_security_group.security_group_id}"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test"
}