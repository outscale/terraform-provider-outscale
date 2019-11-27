#resource "outscale_net" "outscale_net" {
resource "outscale_lin" "outscale_net" {
    count = 1

    ip_range = "10.10.0.0/24"
}

#resource "outscale_net" "outscale_net2" {
resource "outscale_lin" "outscale_net2" {
    count = 1

    ip_range = "10.31.0.0/16"
}

#resource "outscale_net_peering" "outscale_net_peering" {
resource "outscale_lin_peering" "outscale_net_peering" {
    #accepter_net_id   = outscale_net.outscale_net.lin_id
    accepter_net_id   = outscale_lin.outscale_net.lin_id
    #net_id            = outscale_net.outscale_net2.lin_id
    lin_id            = outscale_lin.outscale_net2.lin_id
}

#data "outscale_net_peering" "outscale_net_peering" {
data "outscale_lin_peering" "outscale_net_peering" {
    net_peering_ids = outscale_lin_peering.outscale_net_peering.net_peering_id
}
