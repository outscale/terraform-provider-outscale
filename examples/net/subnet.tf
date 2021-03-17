resource "outscale_subnet" "my_subnet" {
  net_id   = outscale_net.my_net.net_id
  ip_range = var.subnet_ip_range
  subregion_name = "${var.region}a"
}
