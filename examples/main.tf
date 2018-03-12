resource "outscale_inbound_rule" "outscale_inbound_rule1" {
    ip_permissions = {
        from_port = 22
        to_port = 22
        ip_protocol = "tcp"
        ip_ranges = ["46.231.147.8/32"]
    }

    group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule2" {
     ip_permissions = {
        from_port = 443
        to_port = 443
        ip_protocol = "tcp"
        ip_ranges = ["46.231.147.8/32"]
    }

    group_id = "${outscale_firewall_rules_set.outscale_firewall_rules_set.id}"
}

resource "outscale_firewall_rules_set" "outscale_firewall_rules_set" {
    group_description = "test group tf"
    group_name = "sg1-test-group_test1"
}

data "outscale_firewall_rules_set" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.outscale_firewall_rules_set.group_name}"]
		}
	}