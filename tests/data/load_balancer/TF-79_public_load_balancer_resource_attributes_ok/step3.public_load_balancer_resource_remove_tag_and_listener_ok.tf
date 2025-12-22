resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 80
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 80
    }
}
