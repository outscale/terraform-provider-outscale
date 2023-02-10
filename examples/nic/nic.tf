resource "outscale_nic" "eth0" {
  subnet_id          = outscale_subnet.public.subnet_id
  security_group_ids = [outscale_security_group.eth0.id]

  private_ips {
    is_primary = true
    private_ip = "192.168.0.10"
  }
}


resource "outscale_nic" "eth1" {
  subnet_id          = outscale_subnet.private.subnet_id
  security_group_ids = [outscale_security_group.eth1.id]

  private_ips {
    is_primary = true
    private_ip = "192.168.1.201"
  }
}