resource "outscale_ca" "ca_for_aar" {
  ca_pem      = file("certs/certificate.pem")
}

resource "outscale_ca" "ca_for_aar_2" {
  ca_pem      = file("certs/certificate.pem")
}


resource "outscale_api_access_rule" "aar_1" {

  ca_ids = [outscale_ca.ca_for_aar_2.ca_id]

  ip_ranges = ["192.22.0.0/16"]

  cns = ["test-TF201"]

  description = "API Access rules-TF-201-update"

}
