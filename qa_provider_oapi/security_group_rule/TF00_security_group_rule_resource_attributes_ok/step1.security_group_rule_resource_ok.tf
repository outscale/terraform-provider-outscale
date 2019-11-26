resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test"
}

resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    security_group_id = "${outscale_security_group.outscale_security_group.id}"

    from_port_range   = "0"
    to_port_range     = "0"
    #ip_protocol       = "-1"
    ip_protocol       = "tcp"
    ip_range          = "0.0.0.0/0"
}