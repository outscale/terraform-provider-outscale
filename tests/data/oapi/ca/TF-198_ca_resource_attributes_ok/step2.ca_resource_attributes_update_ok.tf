resource "outscale_ca" "test_ca_1" {
  ca_pem      = file("certs/certificate.pem")
  description = "test-ca-update"
}
