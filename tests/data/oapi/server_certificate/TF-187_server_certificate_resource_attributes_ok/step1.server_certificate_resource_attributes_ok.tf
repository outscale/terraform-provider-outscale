resource "outscale_server_certificate" "my_server_certificate" {
   name                   =  "Certificate-TF187"
   body                   =  file("certs/certificate.pem")
   chain                  =  file("certs/certificate.pem")
   private_key            =  file("certs/certificate.key")
}
