resource "outscale_ca" "test_ca_1" {
  ca_pem      = file("data/cert_example/certificate.pem")
  description = "test-ca-update"
}

