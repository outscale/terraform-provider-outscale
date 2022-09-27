resource "outscale_security_group" "my_sg" {
  description = "test security group"
}

resource "outscale_security_group_rule" "my_sg_rule" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.my_sg.id
  rules {
    from_port_range = "22"
    to_port_range   = "22"
    ip_protocol     = "tcp"
    ip_ranges       = ["0.0.0.0/0"]
  }
  rules {
    from_port_range = "80"
    to_port_range   = "80"
    ip_protocol     = "tcp"
    ip_ranges       = ["0.0.0.0/0"]
  }
}
