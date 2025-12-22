resource "outscale_security_group" "outscale_security_group" {
    description         = "test group"
    security_group_name = "sg1-test-group_test"
}
resource "outscale_security_group_rule" "outscale_security_group_rule_https" {
    flow              = "Inbound"
    from_port_range   = 443
    to_port_range     = 443
    ip_protocol       = "tcp"
    ip_range          = "192.168.0.1/32"
    security_group_id = outscale_security_group.outscale_security_group.security_group_id
}
