######Step1:Create the needed resources###
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
    tenancy = "default"
    tags {
     key = "name"
     value = "TF02-NET"
    }
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
    security_group_name = "test-sg-${random_string.suffix[0].result}"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_terraform_lbu_sg"
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


resource "outscale_load_balancer" "internet_facing_lbu_1" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
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
  security_groups = [outscale_security_group.outscale_security_group.id]
  load_balancer_type = "internet-facing"
  tags {
     key = "name"
     value = "lbu-internet-facing-TF02"
    }
 depends_on = [outscale_route.outscale_route,outscale_route_table_link.outscale_route_table_link]
}

resource "outscale_vm" "vm03" {
    image_id           = var.image_id
    vm_type            = var.vm_type
    security_group_ids = [outscale_security_group.outscale_security_group.id]
    subnet_id          = outscale_subnet.subnet-1.subnet_id
}

resource "outscale_public_ip" "public_ip02" {
tags {
     key                    = "name"
     value                  = "EIP-TF02"
  }
}

resource "outscale_public_ip_link" "public_ip_link02" {
    vm_id     = outscale_vm.vm03.vm_id
    public_ip = outscale_public_ip.public_ip02.public_ip
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms02" {
    load_balancer_name = outscale_load_balancer.internet_facing_lbu_1.load_balancer_name
    backend_ips     = [outscale_public_ip.public_ip02.public_ip]
depends_on = [outscale_public_ip_link.public_ip_link02]
}

data "outscale_load_balancer" "load_balancer_TF02" {
    filter {
        name   = "load_balancer_names"
        values = [outscale_load_balancer.internet_facing_lbu_1.load_balancer_name]
    }
depends_on=[outscale_load_balancer_vms.outscale_load_balancer_vms02]
}
