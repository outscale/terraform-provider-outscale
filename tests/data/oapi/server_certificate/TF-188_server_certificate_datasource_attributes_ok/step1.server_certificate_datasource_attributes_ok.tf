resource "outscale_server_certificate" "my_server_certificate2" {
   name                   =  "Certificate-TF188"
   body                   =  file("certs/certificate.pem")
   chain                  =  file("certs/certificate.pem")
   private_key            =  file("certs/certificate.key")
   path                   =  "/terraform/test/"
}

data "outscale_server_certificate" "my_server_certificate" {
      filter {
        name     = "paths"
        values   = [outscale_server_certificate.my_server_certificate2.path]
    }
depends_on = [outscale_server_certificate.my_server_certificate2]
}
