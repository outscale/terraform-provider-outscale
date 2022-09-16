resource "outscale_ca" "ca_for_aar" {
  ca_pem      = file("data/ca/TF-197_ca_datasource_attributes_ok/terraform-ca-certificate.pem.crt")
}

resource "outscale_ca" "ca_for_aar_2" {
  ca_pem      = file("data/ca/TF-197_ca_datasource_attributes_ok/terraform-ca-certificate.pem.crt")
}


resource "outscale_api_access_rule" "aar_1" {

  ca_ids = [outscale_ca.ca_for_aar.ca_id,outscale_ca.ca_for_aar_2.ca_id]

  ip_ranges = ["192.168.2.134", "192.22.0.0/16"]

  cns = ["test-TF201", "test-TF201-2"]

  description = "API Access rules-TF-201"

}


