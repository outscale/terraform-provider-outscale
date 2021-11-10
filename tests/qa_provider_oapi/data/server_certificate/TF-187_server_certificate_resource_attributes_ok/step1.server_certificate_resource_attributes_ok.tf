resource "outscale_server_certificate" "my_server_certificate" { 
   name                   =  "Certificate-TF187"
   body                   =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert.pem")
   chain                  =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-cert-chain.pem")
   private_key            =  file("data/server_certificate/TF-187_server_certificate_resource_attributes_ok/test-key.pem")
}
