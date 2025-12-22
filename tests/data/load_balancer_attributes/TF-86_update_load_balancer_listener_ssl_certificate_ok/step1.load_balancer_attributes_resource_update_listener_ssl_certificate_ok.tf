resource "outscale_server_certificate" "my_server_certificate_TF-86" {
   name                   =  "Certificate-TF86"
   body                   =  file("data/cert_example/certificate.pem")
   private_key            =  file("data/cert_example/certificate.key")
}

resource "outscale_server_certificate" "my_server_certificate_TF-86_2" {
   name                   =  "Certificate-TF86-2"
   body                   =  file("data/cert_example/certificate.pem")
   private_key            =  file("data/cert_example/certificate.key")
}


resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name = "test-lb-${random_string.suffix[0].result}"
   subregion_names= ["${var.region}a"]
listeners {
     backend_port = 8080
     backend_protocol= "HTTPS"
     load_balancer_protocol= "HTTPS"
     load_balancer_port = 8080
     server_certificate_id = outscale_server_certificate.my_server_certificate_TF-86.orn
    }
 tags {
    key = "name"
    value = "public_lbu1"
   }
}

resource "outscale_load_balancer_attributes" "attributes-ssl-certificate" {
   load_balancer_name = outscale_load_balancer.public_lbu1.id
   load_balancer_port = 8080
   server_certificate_id = outscale_server_certificate.my_server_certificate_TF-86_2.orn
}
