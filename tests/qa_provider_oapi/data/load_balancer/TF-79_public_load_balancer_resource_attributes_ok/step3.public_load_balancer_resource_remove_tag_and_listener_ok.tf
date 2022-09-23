resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name ="lbu-TF-79"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 80
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 80
    }
}

