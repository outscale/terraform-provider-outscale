resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name ="lbu-TF-79"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 1024
     backend_protocol= "TCP"
     load_balancer_protocol= "TCP"
     load_balancer_port = 1024
    }
  listeners {
     backend_port = 80
     backend_protocol= "HTTP"
     load_balancer_protocol= "HTTP"
     load_balancer_port = 80
    }
  listeners {
     backend_port            = 8080
     backend_protocol        = "HTTP"
     load_balancer_protocol  = "HTTP"
     load_balancer_port      = 8080
    }

  tags {
     key = "name"
     value = "public_lbu1"
    }
  tags {
     key = "test"
     value = "test-tag"
   }
  tags {
     key = "test-2"
     value = "test-tag-2"
   }
}
