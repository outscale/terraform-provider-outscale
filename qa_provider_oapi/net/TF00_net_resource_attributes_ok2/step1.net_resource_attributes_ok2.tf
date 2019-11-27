resource "outscale_net" "outscale_net" {
    count = 2
    
    ip_range = "10.0.0.0/16"

    tags {
        key   = "Name"
        value = "outscale_net_resource"
    }
}
