resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    net_id     = "${outscale_net.outscale_net.net_id}"
    ip_range = "10.0.0.0/24"
}

resource "outscale_subnet" "outscale_subnet2" {
    net_id     = "${outscale_net.outscale_net.net_id}"
    ip_range = "10.0.1.0/24"
}

resource "outscale_public_ip" "outscale_public_ip" {
}

resource "outscale_public_ip" "outscale_public_ip2" {
}

resource "outscale_nat_service" "outscale_nat_service" {
    depends_on = ["outscale_route.outscale_route"]
    subnet_id     = "${outscale_subnet.outscale_subnet.subnet_id}"
    public_ip_id = "${outscale_public_ip.outscale_public_ip.public_ip_id}"
}

resource "outscale_nat_service" "outscale_nat_service2" {
    depends_on = ["outscale_route.outscale_route2"]
    subnet_id     = "${outscale_subnet.outscale_subnet2.subnet_id}"
    public_ip_id = "${outscale_public_ip.outscale_public_ip2.public_ip_id}"
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_route_table" "outscale_route_table2" {
    net_id = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_route" "outscale_route" {
    destination_ip_range = "0.0.0.0/0"
    gateway_id           = "${outscale_internet_service.outscale_internet_service.internet_service_id}"
    route_table_id       = "${outscale_route_table.outscale_route_table.id}"
}

resource "outscale_route" "outscale_route2" {
    destination_ip_range = "0.0.0.0/0"
    gateway_id           = "${outscale_internet_service.outscale_internet_service.internet_service_id}"
    route_table_id       = "${outscale_route_table.outscale_route_table2.id}"
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    subnet_id      = "${outscale_subnet.outscale_subnet.subnet_id}"
    route_table_id = "${outscale_route_table.outscale_route_table.id}"
}

resource "outscale_route_table_link" "outscale_route_table_link2" {
    subnet_id      = "${outscale_subnet.outscale_subnet2.subnet_id}"
    route_table_id = "${outscale_route_table.outscale_route_table2.id}"
}

resource "outscale_internet_service" "outscale_internet_service" {
}

#resource "outscale_lin_internet_gateway" "outscale_lin_internet_gateway2" {
#}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    net_id              = "${outscale_net.outscale_net.net_id}"
    internet_service_id = "${outscale_internet_service.outscale_internet_service.internet_service_id}"
}

#resource "outscale_lin_internet_gateway_link" "outscale_lin_internet_gateway_link2" {
#    vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
#    internet_gateway_id = "${outscale_lin_internet_gateway.outscale_lin_internet_gateway2.id}"
#}

data "outscale_nat_services" "outscale_nat_services" {
    filter {
        name   = "nat_service_ids"
        values = ["${outscale_nat_service.outscale_nat_service.nat_service_id}", "${outscale_nat_service.outscale_nat_service2.nat_service_id}"]
    }     
}
