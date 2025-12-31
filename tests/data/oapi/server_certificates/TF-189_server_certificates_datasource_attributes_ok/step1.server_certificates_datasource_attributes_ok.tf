resource "outscale_server_certificate" "my_server_certificate3-1" {
   name                   =  "Certificate-TF189-1"
   body                   =  file("certs/certificate.pem")
   chain                  =  file("certs/certificate.pem")
   private_key            =  file("certs/certificate.key")
   path                   =  "/terraform/test1/"
}

resource "outscale_server_certificate" "my_server_certificate3-2" {
   name                   =  "Certificate-TF189-2"
   body                   =  file("certs/certificate.pem")
   chain                  =  file("certs/certificate.pem")
   private_key            =  file("certs/certificate.key")
   path                   =  "/terraform/test2/"
}


data "outscale_server_certificates" "my_server_certificates" {
      filter {
        name     = "paths"
        values   = [outscale_server_certificate.my_server_certificate3-1.path,outscale_server_certificate.my_server_certificate3-2.path]
    }
depends_on = [outscale_server_certificate.my_server_certificate3-1,outscale_server_certificate.my_server_certificate3-2]
}
