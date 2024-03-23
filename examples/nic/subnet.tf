resource "outscale_subnet" "public" {
  net_id         = outscale_net.my_net.net_id
  ip_range       = var.subnet_public_ip_range
  subregion_name = "${var.region}a"
}

resource "outscale_subnet" "private" {
  net_id         = outscale_net.my_net.net_id
  ip_range       = var.subnet_private_ip_range
  subregion_name = "${var.region}a"
}