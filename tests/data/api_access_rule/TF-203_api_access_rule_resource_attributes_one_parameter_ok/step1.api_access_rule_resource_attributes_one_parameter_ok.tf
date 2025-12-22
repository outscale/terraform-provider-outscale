resource "outscale_ca" "new_ca" {
  ca_pem      = file("certs/certificate.pem")
}

resource "outscale_api_access_rule" "aar_1" {

  ca_ids = [outscale_ca.new_ca.ca_id]

  cns = ["test-TF203"]

}
