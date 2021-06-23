resource "outscale_net" "outscale_net" {
    #count = 1

    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    #count = 1

    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.0.0/16"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "terraform-subnet"
    }
    tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_subnet" "outscale_subnet" {
    filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id]
    }
filter {
        name   = "available_ips_counts"
        values = [outscale_subnet.outscale_subnet.available_ips_count]
    }
filter {
        name   = "ip_ranges"
        values = [outscale_subnet.outscale_subnet.ip_range]
    }
filter {
        name   = "net_ids"
        values = [outscale_subnet.outscale_subnet.net_id]
    }
filter {
        name   = "subregion_names"
        values = [outscale_subnet.outscale_subnet.subregion_name]
    }
filter {
        name   = "states"
        values = [outscale_subnet.outscale_subnet.state]
    }

}

