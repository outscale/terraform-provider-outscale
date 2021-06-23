resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name ="lbu-TF-86"
   subregion_names= ["${var.region}a"]
listeners {
     backend_port = 8080
     backend_protocol= "HTTPS"
     load_balancer_protocol= "HTTPS"
     load_balancer_port = 8080
     server_certificate_id = var.server_certificate_id
    }
 tags {
    key = "name"
    value = "public_lbu1"
   }
}

resource "outscale_load_balancer_attributes" "attributes-ssl-certificate" {
   load_balancer_name = outscale_load_balancer.public_lbu1.id
   load_balancer_port = 8080
   server_certificate_id = var.server_certificate_id_2 
}
