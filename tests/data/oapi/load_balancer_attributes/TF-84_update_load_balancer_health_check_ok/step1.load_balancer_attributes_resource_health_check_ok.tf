resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name = "test-lb-${random_string.suffix[0].result}"
   subregion_names= ["${var.region}a"]
   listeners {
     backend_port = 80
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 80
    }
 listeners {
     backend_port = 1024
     backend_protocol= "TCP"
     load_balancer_protocol= "TCP"
     load_balancer_port = 1024
    }
 tags {
    key = "name"
    value = "public_lbu1"
   }
 tags {
    key = "test"
    value = "tags"
   }
}

resource "outscale_load_balancer_attributes" "attributes-health-check" {
   load_balancer_name      = outscale_load_balancer.public_lbu1.id
    health_check  {
        healthy_threshold   = 10
        check_interval      = 36
        path                = "/index.html"
        port                = 80
        protocol            = "HTTP"
        timeout             = 6
        unhealthy_threshold = 6
    }
}
