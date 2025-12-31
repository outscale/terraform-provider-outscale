resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 80
     backend_protocol= "TCP"
     load_balancer_protocol= "TCP"
     load_balancer_port = 80
    }
  tags {
     key = "name"
     value = "public_lbu1"
   }
  tags {
     key   = "test-1"
     value = "outscale_terraform_lbu"
   }
}

resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-1"
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

resource "outscale_load_balancer" "private_lbu_1" {
  load_balancer_name = "test-lb-${random_string.suffix[1].result}"
  listeners {
     backend_port = 8080
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 8080
   }
  subnets = [outscale_subnet.subnet-1.subnet_id]
  security_groups = [outscale_security_group.outscale_security_group.security_group_id]
  load_balancer_type = "internal"
  tags {
     key = "name"
     value = "lbu-internal"
   }
}

data "outscale_load_balancers" "outscale_load_balancers" {
 filter {
        name = "load_balancer_names"
        values = [outscale_load_balancer.public_lbu1.load_balancer_name,outscale_load_balancer.private_lbu_1.load_balancer_name]
      }
depends_on = [outscale_load_balancer.public_lbu1,outscale_load_balancer.private_lbu_1]
}
