resource "outscale_lin" "outscale_lin1" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table1" {
  vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
}

resource "outscale_route_table" "outscale_route_table2" {
  vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
}

data "outscale_route_tables" "outscale_route_tables" {
  route_table_id = ["${outscale_route_table.outscale_route_table1.route_table_id}", "${outscale_route_table.outscale_route_table2.route_table_id}"]
}

resource "outscale_lin" "outscale_lin2" {
  count = 1

  cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet1" {
  count = 1

  availability_zone = "eu-west-2a"
  cidr_block        = "10.0.0.0/16"
  vpc_id            = "${outscale_lin.outscale_lin2.vpc_id}"
}

output "outscale_subnet1" {
  value = "${outscale_subnet.outscale_subnet1.subnet_id}"
}
