resource "outscale_server_certificate" "my_server_certificate" {
  name          = "certificate-${random_string.suffix[0].result}"
  body          =  file("certs/certificate.pem")
  chain         =  file("certs/certificate.pem")
  private_key   =  file("certs/certificate.key")
  path          =  "/terraform/"
}
