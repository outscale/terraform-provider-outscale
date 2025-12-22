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
