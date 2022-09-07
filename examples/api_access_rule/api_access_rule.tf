resource "outscale_api_access_rule" "my_api_rule" {
 # ca_ids      = var.ca_ids
  ip_ranges   = var.ip_ranges
  description = var.description
}
