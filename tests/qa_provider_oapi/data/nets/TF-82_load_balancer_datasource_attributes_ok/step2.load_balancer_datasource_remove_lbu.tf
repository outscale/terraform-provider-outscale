resource "outscale_net" "outscale_net" {
    ip_range = "10.20.0.0/16"
}

resource "outscale_security_group" "outscale_security_group" {
    description         = "test lbu-1"
    security_group_name = "sg1-terraform-lbu-test"
    net_id              = outscale_net.outscale_net.net_id
    tags {
        key   = "Name"
        value = "outscale_terraform_lbu_sg"
    }
}

resource "outscale_subnet" "subnet-1" {
  net_id   = outscale_net.outscale_net.net_id
  ip_range = "10.20.0.0/24"
  tags {
        key   = "Name"
        value = "outscale_terraform_lbu_subnet"
    }
}
