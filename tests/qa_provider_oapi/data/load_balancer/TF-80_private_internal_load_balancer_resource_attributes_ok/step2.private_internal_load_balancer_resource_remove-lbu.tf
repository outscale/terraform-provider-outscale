resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-1"
    security_group_name = "sg1-terraform-lbu-test"
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
   load_balancer_name ="lbu-TF-80"
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
  load_balancer_type = "internal"
  tags {
      key = "name"
      value = "lbu-internal"
   }
}
