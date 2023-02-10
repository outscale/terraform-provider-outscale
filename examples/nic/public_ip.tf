resource "outscale_public_ip" "my_public_ip" {
}

resource "outscale_public_ip_link" "eth0" {
  nic_id = outscale_nic.eth0.nic_id
  public_ip_id = outscale_public_ip.my_public_ip.id
}