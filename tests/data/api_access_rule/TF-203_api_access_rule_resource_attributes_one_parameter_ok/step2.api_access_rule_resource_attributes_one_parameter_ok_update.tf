resource "outscale_ca" "new_ca" {
  ca_pem      = file("data/cert_example/certificate.pem")
}

resource "outscale_api_access_rule" "aar_1" {

ip_ranges = ["172.3.4.0/32"]

}

