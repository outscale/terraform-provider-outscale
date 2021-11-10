resource "outscale_server_certificate" "my_server_certificate2" { 
   name                   =  "Certificate-TF188"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   chain                  =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert-chain.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
   path                   =  "/terraform/test/"
}

data "outscale_server_certificate" "my_server_certificate" { 
      filter {
        name     = "paths"
        values   = [outscale_server_certificate.my_server_certificate2.path]
    }  
depends_on = [outscale_server_certificate.my_server_certificate2]               
}


