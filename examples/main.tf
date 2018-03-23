resource "outscale_outbound_rule" "outscale_outbound_rule1" {
  ip_permissions = {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    ip_ranges   = ["46.231.147.8/32"]
  }

  group_id = "${outscale_firewall_rules_sets.outscale_firewall_rules_sets.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule1" {
  ip_permissions = {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    ip_ranges   = ["46.231.147.8/32"]
  }

  group_id = "${outscale_firewall_rules_sets.outscale_firewall_rules_sets.id}"
}

resource "outscale_inbound_rule" "outscale_inbound_rule2" {
  ip_permissions = {
    from_port   = 443
    to_port     = 443
    ip_protocol = "tcp"
    ip_ranges   = ["46.231.147.8/32"]
  }

  group_id = "${outscale_firewall_rules_sets.outscale_firewall_rules_sets.id}"
}

resource "outscale_firewall_rules_sets" "outscale_firewall_rules_sets" {
  group_description = "Used in the terraform acceptance tests"
  group_name        = "test-1234"
  vpc_id            = "vpc-e9d09d63"

  tags = {
    Name = "tf-acctest"
    Seed = "1234"
  }
}

data "outscale_firewall_rules_sets" "by_filter" {
  filter {
    name   = "group-name"
    values = ["${outscale_firewall_rules_sets.outscale_firewall_rules_sets.group_name}"]
  }
}
