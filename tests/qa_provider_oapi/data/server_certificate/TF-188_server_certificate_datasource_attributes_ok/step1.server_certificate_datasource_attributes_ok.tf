resource "outscale_server_certificate" "my_server_certificate2" { 
   name                   =  "Certificate-TF188"
   body                   =  file("data/cert_example/certificate.pem")
   chain                  =  file("data/cert_example/certificate.pem")
   private_key            =  file("data/cert_example/certificate.key")
   path                   =  "/terraform/test/"
}

data "outscale_server_certificate" "my_server_certificate" { 
      filter {
        name     = "paths"
        values   = [outscale_server_certificate.my_server_certificate2.path]
    }  
depends_on = [outscale_server_certificate.my_server_certificate2]               
}


