resource "outscale_server_certificate" "new_server_certificate" {
   name                   =  "Certificate-TF-79"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
}

resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name ="lbu-TF-79"
  subregion_names= ["${var.region}a"]
  listeners {
     backend_port = 8080
     backend_protocol= "HTTPS"
     load_balancer_protocol= "HTTPS"
     load_balancer_port = 8080
     server_certificate_id = outscale_server_certificate.my_server_certificate.orn
    }
}
