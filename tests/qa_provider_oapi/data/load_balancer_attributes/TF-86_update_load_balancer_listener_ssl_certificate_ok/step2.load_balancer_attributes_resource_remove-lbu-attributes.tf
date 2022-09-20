resource "outscale_server_certificate" "server_certificate_1" {
   name                   =  "Certificate-TF86-1"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
}

resource "outscale_server_certificate" "server_certificate_2" {
   name                   =  "Certificate-TF86-2"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
}
resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name ="lbu-TF-86"
   subregion_names= ["${var.region}a"]
listeners {
     backend_port = 8080
     backend_protocol= "HTTPS"
     load_balancer_protocol= "HTTPS"
     load_balancer_port = 8080
     server_certificate_id = outscale_server_certificate.server_certificate_2.orn
    }
 tags {
    key = "name"
    value = "public_lbu1"
   }
}
