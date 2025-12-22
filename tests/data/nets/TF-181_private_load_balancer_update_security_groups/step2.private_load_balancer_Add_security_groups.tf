resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-TF181"
    security_group_name = "terraform-sg-lbu-TF181-1"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_security_group" "outscale_security_group-2" {
    description         = "test lbu-2"
    security_group_name = "terraform-sg-lbu-TF181-2"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "subnet-1" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.0.0.0/24"
}
resource "outscale_load_balancer" "private_lbu" {
  load_balancer_name ="lbu-TF-181-${var.suffixe_lbu_name}"
  listeners {
     backend_port = 80
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 80
    }
  subnets = [outscale_subnet.subnet-1.subnet_id]
 security_groups = [outscale_security_group.outscale_security_group.security_group_id,outscale_security_group.outscale_security_group-2.security_group_id]
  load_balancer_type = "internal"
  tags {
     key = "name"
     value = "lbu-update-sg"
    }
}

