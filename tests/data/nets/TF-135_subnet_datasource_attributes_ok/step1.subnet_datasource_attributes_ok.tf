resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {

    subregion_name = "${var.region}a"
    ip_range       = "10.0.1.0/24"
    net_id         = outscale_net.outscale_net.net_id
    tags {
     key = "name"
     value = "terraform-subnet"
    }
    tags {
      key = "tag:key"
      value = "value:tags"
     }
}

resource "outscale_subnet" "outscale_subnet-2" {

    subregion_name = "${var.region}a"
    ip_range       = "10.0.2.0/28"
    net_id         = outscale_net.outscale_net.net_id
}


data "outscale_subnet" "outscale_subnet" {
    filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id]
    }
}

data "outscale_subnet" "outscale_subnet-2" {
filter {
        name   = "available_ips_counts"
        values = [outscale_subnet.outscale_subnet-2.available_ips_count]
    }
filter {
        name   = "net_ids"
        values = [outscale_subnet.outscale_subnet.net_id]
    }
depends_on = [outscale_subnet.outscale_subnet-2]
}

data "outscale_subnet" "outscale_subnet-3" {
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

data "outscale_subnet" "outscale_subnet-4" {
filter {
        name   = "subregion_names"
        values = [outscale_subnet.outscale_subnet-2.subregion_name]
    }
filter {
        name   = "states"
        values = [outscale_subnet.outscale_subnet.state]
    }
filter {
        name   = "net_ids"
        values = [outscale_subnet.outscale_subnet.net_id]
    }
filter {
        name   = "subnet_ids"
        values = [outscale_subnet.outscale_subnet.subnet_id]
    }
depends_on = [outscale_subnet.outscale_subnet-2]
}


data "outscale_subnet" "outscale_subnet-5" {
filter {
        name   = "tag_keys"
        values = ["tag:key"]
    }
depends_on = [outscale_subnet.outscale_subnet]
}

data "outscale_subnet" "outscale_subnet-6" {
filter {
        name   = "tag_values"
        values = ["value:tags"]
    }
depends_on = [outscale_subnet.outscale_subnet]
}

data "outscale_subnet" "outscale_subnet-7" {
filter {
        name   = "tags"
        values = ["tag:key=value:tags"]
    }
depends_on = [outscale_subnet.outscale_subnet]
}
