resource "outscale_nat_service" "outscale_nat_service" {
    count = 1

    subnet_id = "subnet-fake"
}

output "outscale_nat_service" {
  value = outscale_nat_service.outscale_nat_service.allocation_id
}
