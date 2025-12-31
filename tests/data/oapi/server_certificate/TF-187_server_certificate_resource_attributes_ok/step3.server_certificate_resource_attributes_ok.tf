resource "outscale_server_certificate" "my_server_certificate" {
   name                   =  "Cert-TF187-update"
   body                   =  file("certs/certificate.pem")
   chain                  =  file("certs/certificate.pem")
   private_key            =  file("certs/certificate.key")
   path                   =  "/terraform/"
}
