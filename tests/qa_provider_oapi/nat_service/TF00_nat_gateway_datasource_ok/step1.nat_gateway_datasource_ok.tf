resource "outscale_net" "outscale_net" {
  count = 1
  
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  count = 1

  net_id = outscale_net.outscale_net.vpc_id
}

resource "outscale_nat_service" "outscale_nat_service" {
  count = 1

  subnet_id = outscale_subnet.outscale_subnet.subnet_id
}


output "nat_gateway" {
  value = outscale_nat_service.outscale_nat_service.nat_gateway_id
}
