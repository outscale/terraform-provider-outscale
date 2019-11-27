resource "outscale_lin" "outscale_lin" {
    count = 1
  
    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
  count = 1

  vpc_id = outscale_lin.outscale_lin.vpc_id
}

resource "outscale_public_ip" "outscale_public_ip" {
    count = 1

    domain = ""
}

resource "outscale_nat_service" "outscale_nat_service" {
    count = 1

    #allocation_id = ""
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

output "outscale_nat_service" {
    value = outscale_nat_service.outscale_nat_service.allocation_id
}
