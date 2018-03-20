# resource "outscale_lin" "vpc" {
#   cidr_block = "10.0.0.0/16"
# }

resource "outscale_lin_attributes" "outscale_lin_attributes" {
  attribute            = "enableDnsSupport"
  enable_dns_hostnames = true
  vpc_id               = "vpc-e9d09d63"
}
