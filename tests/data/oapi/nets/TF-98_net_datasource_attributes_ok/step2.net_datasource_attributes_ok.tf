resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/18"
    
    tags {
        key = "Name"
        value = "outscale_net_resource"
    }

    tags {
        key = "VerSion"
        value = "1Q84"
    }
}

data "outscale_net" "outscale_net" {
    filter {
        name   = "net_ids"
        values = [outscale_net.outscale_net.net_id]
    }
}

data "outscale_net" "outscale_net_2" {
    filter {
        name   = "tags"
        values = ["VerSion=1Q84"]
    }
}

data "outscale_net" "outscale_net_3" {
    filter {
        name   = "tag_keys"
        values = ["VerSion"]
    }
}

data "outscale_net" "outscale_net_4" {
    filter {
        name   = "tag_values"
        values = ["1Q84"]
    }
}
data "outscale_net" "outscale_net_5" {
    filter {
        name   = "ip_ranges"
        values = ["10.0.0.0/18"]
    }
}
