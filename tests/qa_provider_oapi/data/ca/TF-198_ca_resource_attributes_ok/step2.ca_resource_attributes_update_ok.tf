resource "outscale_ca" "test_ca_1" {
  ca_pem      = file("data/ca/TF-197_ca_datasource_attributes_ok/terraform-ca-certificate.pem.crt")
  description = "test-ca-update"
}

