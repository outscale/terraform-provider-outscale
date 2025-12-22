resource "outscale_server_certificate" "my_server_certificate" { 
   name                   =  "Cert-TF187-update"
   body                   =  file("data/cert_example/certificate.pem")
   chain                  =  file("data/cert_example/certificate.pem")
   private_key            =  file("data/cert_example/certificate.key")
   path                   =  "/terraform/"
}
