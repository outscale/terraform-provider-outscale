resource "outscale_net" "outscale_net" {
    ip_range = "10.10.0.0/24"
}

resource "outscale_net" "outscale_net2" {
    ip_range = "10.32.0.0/16"
}

resource "outscale_net_peering" "outscale_net_peering" {
    accepter_net_id   = outscale_net.outscale_net.net_id
    source_net_id            = outscale_net.outscale_net2.net_id
}
