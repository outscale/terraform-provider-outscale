resource "outscale_net" "outscale_net_sg" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_security_group" "private_sg" {
    remove_default_outbound_rule = true
    description         = "test group-TF204"
    security_group_name = "test-sg-${random_string.suffix[0].result}"
    net_id              = outscale_net.outscale_net_sg.net_id
}
