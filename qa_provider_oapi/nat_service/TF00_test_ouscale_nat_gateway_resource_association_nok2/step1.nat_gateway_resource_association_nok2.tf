resource "outscale_net" "outscale_net" {
  count = 1
  
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  count = 1

  vpc_id = "${outscale_net.outscale_net.vpc_id}"
}

resource "outscale_nat_service" "outscale_nat_service" {
  count = 1

  subnet_id = "${outscale_vm.outscale_vm.subnet_id}"
  public_ip = "171.33.100.89"
}

output "outscale_nat_service" {
  value = "${outscale_nat_service.outscale_nat_service.allocation_id}"
}