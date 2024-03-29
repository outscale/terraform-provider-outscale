resource "outscale_security_group" "my_sg" {
  description = "test security group"
}

resource "outscale_security_group_rule" "my_sg_rule" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.my_sg.id
  rules {
    from_port_range = "3389"
    to_port_range   = "3389"
    ip_protocol     = "tcp"
    ip_ranges       = var.allowed_cidr
  }
}
