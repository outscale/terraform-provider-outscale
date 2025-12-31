resource "outscale_ca" "ca_for_api_access_rule" {
  ca_pem      = file("certs/certificate.pem")
}

resource "outscale_ca" "ca_for_api_access_rule_2" {
  ca_pem      = file("certs/certificate.pem")
}


resource "outscale_api_access_rule" "api_access_rule_1" {

  ca_ids = [outscale_ca.ca_for_api_access_rule.ca_id,outscale_ca.ca_for_api_access_rule_2.ca_id]

  ip_ranges = ["192.168.2.34", "192.14.0.0/16"]

  cns = ["outscale-1", "test-TF200", "test-TF200-2"]

  description = "API Access rules-TF-200"

}


data "outscale_api_access_rule" "data_api_access_rule_1" {

  filter {
    name   = "api_access_rule_ids"
    values = [outscale_api_access_rule.api_access_rule_1.api_access_rule_id]
  }

}


data "outscale_api_access_rule" "data_api_access_rule_2" {

  filter {
    name   = "ca_ids"
    values = [outscale_ca.ca_for_api_access_rule_2.ca_id]
  }

depends_on =[outscale_api_access_rule.api_access_rule_1]
}


data "outscale_api_access_rule" "data_api_access_rule_3" {

  filter {
    name   = "ip_ranges"
    values = ["192.14.0.0/16"]
  }

depends_on =[outscale_api_access_rule.api_access_rule_1]

}

data "outscale_api_access_rule" "data_api_access_rule_4" {

  filter {
    name   = "cns"
    values = ["test-TF200-2"]
  }

depends_on =[outscale_api_access_rule.api_access_rule_1]
}

data "outscale_api_access_rule" "data_api_access_rule_5" {

  filter {
    name   = "descriptions"
    values = ["API Access rules-TF-200"]
  }

depends_on =[outscale_api_access_rule.api_access_rule_1]
}
