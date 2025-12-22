resource "outscale_ca" "new_ca_for_data" {
  ca_pem      = file("data/cert_example/certificate.pem")
  description = "test-TF199"
}

resource "outscale_ca" "new_ca_for_data_2" {
  ca_pem      = file("data/cert_example/certificate.pem")
  description = "test-TF199-2"
}

data "outscale_cas" "data_cas_1" {

  filter {
    name   = "ca_ids"
    values = [outscale_ca.new_ca_for_data.ca_id,outscale_ca.new_ca_for_data_2.ca_id]
  }

}

data "outscale_cas" "data_cas_2" {

  filter {
    name   = "descriptions"
    values = [outscale_ca.new_ca_for_data.description, outscale_ca.new_ca_for_data_2.description]
  }

}
