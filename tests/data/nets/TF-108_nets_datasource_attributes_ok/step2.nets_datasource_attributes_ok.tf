resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/17"

    tags  {
        key   = "Name"
        value = "outscale_net_resource"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_net" "outscale_net2" {
    ip_range = "10.2.0.0/17"

    tags  {
        key   = "Name-net-2"
        value = "outscale_net_resource2"
    }
}

data "outscale_nets" "outscale_nets" {
    filter {
        name   = "net_ids"
        values = [outscale_net.outscale_net.net_id, outscale_net.outscale_net2.net_id]
    }
}

data "outscale_nets" "outscale_nets_2" {    
     filter {
        name   = "tags"
        values = ["Name-net-2=outscale_net_resource2"]
    }
}

data "outscale_nets" "outscale_nets_3" {
    filter {
        name   = "tag_keys"
        values = ["Name-net-2"]
    }
}

data "outscale_nets" "outscale_nets_4" {
    filter {
        name   = "tag_values"
        values = ["outscale_net_resource2"]
    }
}
data "outscale_nets" "outscale_nets_5" {
    filter {
        name   = "ip_ranges"
        values = ["10.2.0.0/17"]
    }
}
