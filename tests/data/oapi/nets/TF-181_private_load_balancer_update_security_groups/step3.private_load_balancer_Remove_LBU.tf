resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-TF181"
    security_group_name = "test-sg-${random_string.suffix[0].result}"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_security_group" "outscale_security_group-2" {
    description         = "test lbu-2"
    security_group_name = "test-sg-${random_string.suffix[1].result}"
    net_id              = outscale_net.outscale_net.net_id
}

resource "outscale_subnet" "subnet-1" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.0.0.0/24"
}
