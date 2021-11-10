resource "outscale_server_certificate" "my_server_certificate3-1" { 
   name                   =  "Certificate-TF189-1"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   chain                  =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert-chain.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
   path                   =  "/terraform/test1/"
}

resource "outscale_server_certificate" "my_server_certificate3-2" {
   name                   =  "Certificate-TF189-2"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   chain                  =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert-chain.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
   path                   =  "/terraform/test2/"
}


data "outscale_server_certificates" "my_server_certificates" { 
      filter {
        name     = "paths"
        values   = [outscale_server_certificate.my_server_certificate3-1.path,outscale_server_certificate.my_server_certificate3-2.path]
    }  
depends_on = [outscale_server_certificate.my_server_certificate3-1,outscale_server_certificate.my_server_certificate3-2]               
}
