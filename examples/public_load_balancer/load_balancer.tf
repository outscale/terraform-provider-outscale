resource "outscale_load_balancer" "my_public_lb" {
  subregion_names    = ["${var.region}a"]
  load_balancer_name = "example-lb-${random_string.suffix.result}"
  listeners {
    backend_port           = 80
    backend_protocol       = "HTTP"
    load_balancer_protocol = "HTTP"
    load_balancer_port     = 80
  }
}

output "load_balancer_url" {
  value = "http://${outscale_load_balancer.my_public_lb.dns_name}"
}

resource "outscale_load_balancer_vms" "backend_vms" {
  load_balancer_name = outscale_load_balancer.my_public_lb.load_balancer_name
  backend_vm_ids     = [for _, vm in outscale_vm.my_vms : vm.vm_id]
}

resource "outscale_load_balancer_attributes" "my_health_check" {
  load_balancer_name = outscale_load_balancer.my_public_lb.load_balancer_name
  health_check {
    healthy_threshold   = 10
    check_interval      = 30
    path                = "/"
    port                = 80
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 5
  }
}
