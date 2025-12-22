resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
    tenancy = "default"
}

resource "outscale_route_table" "outscale_route_table" {
  net_id = outscale_net.outscale_net.net_id
  tags {
     key = "name"
     value = "terraform-RT-lbu"
    }
}
resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-2"
    security_group_name = "TF81-lbu-1"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_terraform_lbu_sg"
    }
}

resource "outscale_security_group" "outscale_security_group-2" {
    description         = "test lbu-2"
    security_group_name = "TF81-lbu-2"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_terraform_lbu_sg_2"
    }
}
resource "outscale_subnet" "subnet-1" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.0.0.0/24"
  tags {
        key   = "Name"
        value = "outscale_terraform_lbu_subnet"
    }
}


resource "outscale_internet_service" "outscale_internet_service" {
}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
    net_id = outscale_net.outscale_net.net_id
}

resource "outscale_route" "outscale_route" {
    gateway_id           = outscale_internet_service.outscale_internet_service.id
    destination_ip_range = "0.0.0.0/0"
   route_table_id       = outscale_route_table.outscale_route_table.route_table_id
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    route_table_id  = outscale_route_table.outscale_route_table.route_table_id
    subnet_id       = outscale_subnet.subnet-1.subnet_id
}


resource "outscale_load_balancer" "private_lbu_1" {
  load_balancer_name ="lbu-TF-81-${var.suffixe_lbu_name}"
  listeners {
     backend_port = 80
     backend_protocol= "TCP"
     load_balancer_protocol= "TCP"
     load_balancer_port = 80
    }
  listeners {
     backend_port = 8080
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 8080
    }
  subnets = [outscale_subnet.subnet-1.subnet_id]
  security_groups = [outscale_security_group.outscale_security_group.id, outscale_security_group.outscale_security_group-2.id]
  load_balancer_type = "internet-facing"
  tags {
     key = "name"
     value = "lbu-internet-facing"
    }
 depends_on = [outscale_route.outscale_route,outscale_route_table_link.outscale_route_table_link]
}
