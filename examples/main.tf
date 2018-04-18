resource "outscale_lin" "outscale_lin1" {
  cidr_block = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
  vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
}

resource "outscale_route_table" "outscale_route_table2" {
  vpc_id = "${outscale_lin.outscale_lin1.vpc_id}"
}

data "outscale_route_tables" "outscale_route_tables" {
  route_table_id = ["${outscale_route_table.outscale_route_table.route_table_id}", "${outscale_route_table.outscale_route_table2.route_table_id}"]
}
