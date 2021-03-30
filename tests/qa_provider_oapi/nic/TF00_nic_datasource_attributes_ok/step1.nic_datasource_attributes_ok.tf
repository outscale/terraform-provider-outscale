resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.0.0/16"
    net_id         = outscale_net.outscale_net.net_id
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = outscale_subnet.outscale_subnet.subnet_id
}

data "outscale_nic" "outscale_nic" {
    filter {
        name = "nic_ids"
        values = [outscale_nic.outscale_nic.nic_id]
    }    
}
