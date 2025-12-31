resource "outscale_ca" "new_ca" {
  ca_pem      = file("certs/certificate.pem")
}

resource "outscale_api_access_rule" "aar_1" {

ip_ranges = ["172.3.4.0/32"]

}
