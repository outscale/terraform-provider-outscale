resource "outscale_net" "outscale_net" {
  
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {

    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.10.0/24"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "subnet-tags-1"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_subnet" "outscale_subnet2" {
  
    subregion_name = format("%s%s", var.region, "a")
    ip_range       = "10.0.20.0/24"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "subnet-tags-2"
    }
}
data "outscale_subnets" "outscale_subnets" {
    filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id, outscale_subnet.outscale_subnet2.subnet_id]
    }
}
data "outscale_subnets" "outscale_subnets-2" {
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
