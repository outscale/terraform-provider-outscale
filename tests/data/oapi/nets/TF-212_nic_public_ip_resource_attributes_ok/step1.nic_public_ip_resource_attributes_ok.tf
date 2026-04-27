resource "outscale_net" "outscale_net" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  subregion_name = "${var.region}a"
  ip_range       = "10.0.0.0/16"
  net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
  subnet_id = outscale_subnet.outscale_subnet.subnet_id
  private_ips {
    is_primary = true
    private_ip = "10.0.67.45"
  }
}

resource "outscale_public_ip" "public_ip01" {
}

resource "outscale_public_ip_link" "public_ip_link01" {
  nic_id    = outscale_nic.outscale_nic.nic_id
  public_ip = outscale_public_ip.public_ip01.public_ip
}
