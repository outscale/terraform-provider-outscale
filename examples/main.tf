resource "outscale_lin_internet_gateway" "res" {}

data "outscale_lin_internet_gateway" "data" {
  internet_gateway_id = "${outscale_lin_internet_gateway.res.internet_gateway_id}"
}
