resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name ="lbu-TF-83"
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
    key = "platform"
    value = "${var.region}a"
   }
}
