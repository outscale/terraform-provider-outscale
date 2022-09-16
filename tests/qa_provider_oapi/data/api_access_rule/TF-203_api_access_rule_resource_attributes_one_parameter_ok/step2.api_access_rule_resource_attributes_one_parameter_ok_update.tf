resource "outscale_ca" "new_ca" {
  ca_pem      = file("data/ca/TF-197_ca_datasource_attributes_ok/terraform-ca-certificate.pem.crt")
}

resource "outscale_api_access_rule" "aar_1" {

ip_ranges = ["172.3.4.0/32"]

}

