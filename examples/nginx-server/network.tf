# Security group attached to the VM.
resource "outscale_security_group" "my-security-group" {
  security_group_name = "${local.name_prefix}-sg"
  description         = "Security group for the Nginx server example"
}

# Allow SSH access.
resource "outscale_security_group_rule" "ssh_inbound" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.my-security-group.security_group_id
  ip_protocol       = "tcp"
  from_port_range   = 22
  to_port_range     = 22
  ip_range          = "0.0.0.0/0"
}

# Allow HTTP access.
resource "outscale_security_group_rule" "http_inbound" {
  flow              = "Inbound"
  security_group_id = outscale_security_group.my-security-group.security_group_id
  ip_protocol       = "tcp"
  from_port_range   = 80
  to_port_range     = 80
  ip_range          = "0.0.0.0/0"
}