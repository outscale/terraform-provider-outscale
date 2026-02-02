resource "outscale_server_certificate" "my_server_certificate_TF-86" {
  name                   = "certificate-${random_string.suffix[0].result}"
  body                   =  file("certs/certificate.pem")
  private_key            =  file("certs/certificate.key")
}

resource "outscale_server_certificate" "my_server_certificate_TF-86_2" {
  name                   = "certificate-${random_string.suffix[1].result}"
  body                   =  file("certs/certificate.pem")
  private_key            =  file("certs/certificate.key")
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
