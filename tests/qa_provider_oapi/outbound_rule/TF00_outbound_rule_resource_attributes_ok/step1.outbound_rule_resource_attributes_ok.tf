resource "outscale_net" "net" {
    ip_range = "10.0.0.0/16"

    tag {
        Name = "outbound rule test"
    }
}

resource "outscale_outbound_rule" "outscale_outbound_rule" {
    from_port_range       = 0
    to_port_range         = 0
    ip_protocol           = "-1"
    ip_ranges             = ["1.2.3.4/32"]
    firewall_rules_set_id = outscale_firewall_rules_set.outscale_firewall_rules_set.id
}

resource "outscale_firewall_rules_set" "outscale_firewall_rules_set" {
    description = "test group"
    name        = "sg1-test-group_test"
    net_id      = outscale_lin.lin.id
}
