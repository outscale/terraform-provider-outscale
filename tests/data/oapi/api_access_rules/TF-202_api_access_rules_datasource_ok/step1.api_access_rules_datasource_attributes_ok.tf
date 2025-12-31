resource "outscale_ca" "ca_for_api_access_rule" {
  ca_pem      = file("certs/certificate.pem")
}

resource "outscale_ca" "ca_for_api_access_rule_2" {
  ca_pem      = file("certs/certificate.pem")
}


resource "outscale_api_access_rule" "aar_1" {

  ca_ids = [outscale_ca.ca_for_api_access_rule.ca_id,outscale_ca.ca_for_api_access_rule_2.ca_id]

  ip_ranges = ["192.168.2.34", "192.14.0.0/16"]

  cns = ["outscale-1", "test-TF202", "test-TF202-2"]

  description = "API Access rules-TF-202-1"

}

resource "outscale_api_access_rule" "aar_2" {

  ca_ids = [outscale_ca.ca_for_api_access_rule.ca_id,outscale_ca.ca_for_api_access_rule_2.ca_id]

  ip_ranges = ["192.168.2.34", "192.14.0.0/16"]

  cns = ["outscale-1", "test-TF202"]

  description = "AAR-TF-202-1"

}

data "outscale_api_access_rules" "data_aar_1" {

  filter {
    name   = "api_access_rule_ids"
    values = [outscale_api_access_rule.aar_1.api_access_rule_id, outscale_api_access_rule.aar_2.api_access_rule_id]
  }

}


data "outscale_api_access_rules" "data_api_access_rule_2" {

  filter {
    name   = "ca_ids"
    values = [outscale_ca.ca_for_api_access_rule_2.ca_id, outscale_ca.ca_for_api_access_rule.ca_id]
  }

depends_on =[outscale_api_access_rule.aar_1,outscale_api_access_rule.aar_2]
}


data "outscale_api_access_rules" "data_api_access_rule_3" {

  filter {
    name   = "ip_ranges"
    values = ["192.14.0.0/16"]
  }

depends_on =[outscale_api_access_rule.aar_1, outscale_api_access_rule.aar_2]

}

data "outscale_api_access_rules" "data_api_access_rule_4" {

  filter {
    name   = "cns"
    values = ["outscale-1"]
  }

depends_on =[outscale_api_access_rule.aar_1, outscale_api_access_rule.aar_2]
}

data "outscale_api_access_rules" "data_api_access_rule_5" {

  filter {
    name   = "descriptions"
    values = ["API Access rules-TF-202", "AAR-TF-202-1"]
  }

depends_on =[outscale_api_access_rule.aar_1, outscale_api_access_rule.aar_2]
}
