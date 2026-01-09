resource "outscale_ca" "new_ca" {
  ca_pem      = file("certs/certificate.pem")
  description = "test-TF197"
}

data "outscale_ca" "data_ca_1" {

  filter {
    name   = "ca_ids"
    values = [outscale_ca.new_ca.ca_id]
  }

}

data "outscale_ca" "data_ca_2" {

  filter {
    name   = "descriptions"
    values = [outscale_ca.new_ca.description]
  }

}
