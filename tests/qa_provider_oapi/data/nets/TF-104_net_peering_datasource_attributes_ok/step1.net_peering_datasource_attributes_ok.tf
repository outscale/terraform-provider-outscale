resource "outscale_net" "outscale_net" {
    #count = 1

    ip_range = "10.10.0.0/24"
}

resource "outscale_net" "outscale_net2" {
    #count = 1

    ip_range = "10.31.0.0/16"
}

resource "outscale_net_peering" "outscale_net_peering" {
    accepter_net_id   = outscale_net.outscale_net.net_id
    source_net_id     = outscale_net.outscale_net2.net_id
    accepter_owner_id = "${var.account_id}"
    tags {
     key = "name"
     value = "test-net-peering"
    }
tags {
      key = "Key-2"
      value = "value-tags-2"
     }
}

data "outscale_net_peering" "outscale_net_peering" {
    filter {
        name = "net_peering_ids"
        values = [outscale_net_peering.outscale_net_peering.net_peering_id]
    }
}
