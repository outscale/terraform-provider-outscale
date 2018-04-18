# resource "outscale_lin" "outscale_lin1" {
#   cidr_block = "10.0.0.0/16"
# }

# resource "outscale_route_table" "outscale_route_table1" {
#   vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
# }

# resource "outscale_route_table" "outscale_route_table2" {
#   vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
# }

# data "outscale_route_tables" "outscale_route_tables" {
#   route_table_id = ["${outscale_route_table.outscale_route_table1.route_table_id}", "${outscale_route_table.outscale_route_table2.route_table_id}"]
# }

# resource "outscale_lin" "outscale_lin2" {
#   count = 1

#   cidr_block = "10.0.0.0/16"
# }

# resource "outscale_subnet" "outscale_subnet1" {
#   count = 1

#   availability_zone = "eu-west-2a"
#   cidr_block        = "10.0.0.0/16"
#   vpc_id            = "${outscale_lin.outscale_lin2.vpc_id}"
# }

# output "outscale_subnet1" {
#   value = "${outscale_subnet.outscale_subnet1.subnet_id}"
# }

resource "outscale_lin" "outscale_lin" {
  cidr_block = "10.0.0.0/24"
}

resource "outscale_route_table" "outscale_route_table" {
  vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway" {}

resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link" {
  internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.internet_gateway_id}"
  vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_route" "outscale_route" {
  gateway_id             = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway.internet_gateway_id}"
  destination_cidr_block = "10.0.0.0/16"
  route_table_id         = "${outscale_route_table.outscale_route_table.route_table_id}"
}
