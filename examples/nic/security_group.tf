resource "outscale_security_group" "eth0" {
  description = "Security Group for public subnet"
  net_id = outscale_net.my_net.id
}

resource "outscale_security_group_rule" "eth0_rule" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.eth0.id
  rules {
    from_port_range = "22"
    to_port_range   = "22"
    ip_protocol     = "tcp"
    ip_ranges       = var.allowed_cidr
  }
}

resource "outscale_security_group" "eth1" {
  description = "Security Group for private subnet"
  net_id = outscale_net.my_net.id
}

resource "outscale_security_group_rule" "eth1_rule" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.eth1.id
  rules {
    from_port_range = "22"
    to_port_range   = "22"
    ip_protocol     = "tcp"
    ip_ranges       = [var.customer_ip_range]
  }
}
