resource "outscale_net" "outscale_net" {
  
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {

    subregion_name = "${var.region}a"
    ip_range       = "10.0.10.0/24"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name-tf137"
     value = "subnet-tags-1"
    }
tags {
      key = "Key"
      value = "value-tags"
     }
}

resource "outscale_subnet" "outscale_subnet2" {
  
    subregion_name = "${var.region}b"
    ip_range       = "10.0.20.0/28"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name-tf137"
     value = "subnet-tags-2"
    }
    tags {
      key = "Key"
      value = "value-tags"
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
        name   = "available_ips_counts"
        values = [outscale_subnet.outscale_subnet2.available_ips_count]
    }
depends_on = [outscale_subnet.outscale_subnet2]
}

data "outscale_subnets" "outscale_subnets-3" {

filter {
        name   = "ip_ranges"
        values = [outscale_subnet.outscale_subnet.ip_range]
    }
filter {
        name   = "net_ids"
        values = [outscale_subnet.outscale_subnet.net_id]
}
depends_on = [outscale_subnet.outscale_subnet]
}

data "outscale_subnets" "outscale_subnets-4" {
 filter {
        name   = "subregion_names"
        values = [outscale_subnet.outscale_subnet.subregion_name]
    }
  filter {
        name   = "net_ids"
        values = [outscale_subnet.outscale_subnet.net_id]
  }
 filter {
        name   = "states"
        values = [outscale_subnet.outscale_subnet.state]
    }
depends_on = [outscale_subnet.outscale_subnet]
}


data "outscale_subnets" "outscale_subnets-5" {
 filter {
        name   = "tag_keys"
        values = ["name-tf137"]
    }
depends_on = [outscale_subnet.outscale_subnet,outscale_subnet.outscale_subnet2]
}

data "outscale_subnets" "outscale_subnets-6" {
 filter {
        name   = "tag_values"
        values = ["subnet-tags-2","subnet-tags-2"]
    }
depends_on = [outscale_subnet.outscale_subnet,outscale_subnet.outscale_subnet2]
}

data "outscale_subnets" "outscale_subnets-7" {
 filter {
        name   = "tags"
        values = ["name-tf137=subnet-tags-1"]
    }
depends_on = [outscale_subnet.outscale_subnet]
}
