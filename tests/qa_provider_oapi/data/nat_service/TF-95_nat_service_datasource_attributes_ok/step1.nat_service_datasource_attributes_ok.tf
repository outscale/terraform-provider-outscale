resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
    tags {
     key = "name"
     value = "net"
    }
}

resource "outscale_subnet" "outscale_subnet" {
    net_id     = outscale_net.outscale_net.net_id
    ip_range = "10.0.0.0/18"
    tags {
     key = "name"
     value = "subnet"
    }
}

resource "outscale_public_ip" "outscale_public_ip" {
 tags {
    key = "name"
    value = "public_ip"
    }
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
     tags {
    key = "name"
    value = "route_table"
    }
}

resource "outscale_route" "outscale_route" {
    destination_ip_range = "0.0.0.0/0"
    gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
    route_table_id       = outscale_route_table.outscale_route_table.route_table_id
}

resource "outscale_route_table_link" "outscale_route_table_link" {
    subnet_id      = outscale_subnet.outscale_subnet.subnet_id
    route_table_id = outscale_route_table.outscale_route_table.id
}

resource "outscale_internet_service" "outscale_internet_service" {
  tags {
    key = "name"
    value = "internet_service"
    }
}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    net_id              = outscale_net.outscale_net.net_id
    internet_service_id = outscale_internet_service.outscale_internet_service.id
}

resource "outscale_nat_service" "outscale_nat_service" {
    depends_on   = [outscale_route.outscale_route]
    subnet_id    = outscale_subnet.outscale_subnet.subnet_id
    public_ip_id = outscale_public_ip.outscale_public_ip.public_ip_id
   tags {
     key = "name"
     value = "nat"
    }
   tags {
      key = "Key"
      value = "value-tags"
     }
}

data "outscale_nat_service" "outscale_nat_service" {
    filter {
        name   = "subnet_ids"
        values = [outscale_nat_service.outscale_nat_service.subnet_id]
    }    
depends_on=[outscale_nat_service.outscale_nat_service]
}

data "outscale_nat_service" "outscale_nat_service-2" {
    filter {
        name   = "nat_service_ids"
        values = [outscale_nat_service.outscale_nat_service.nat_service_id]
    }
}

data "outscale_nat_service" "outscale_nat_service-3" {
    filter {
        name   = "net_ids"
        values = [outscale_nat_service.outscale_nat_service.net_id]
    }
depends_on=[outscale_nat_service.outscale_nat_service]
}

data "outscale_nat_service" "outscale_nat_service-4" {
    filter {
        name   = "states"
        values = [outscale_nat_service.outscale_nat_service.state]
    }
    filter {
        name   = "tag_keys"
        values = ["Key"]
   }
}

data "outscale_nat_service" "outscale_nat_service-5" {
    filter {
        name   = "states"
        values = [outscale_nat_service.outscale_nat_service.state]
    }
    filter {
        name   = "tag_values"
        values = ["value-tags"]
   }
}

data "outscale_nat_service" "outscale_nat_service-6" {
    filter {
        name   = "states"
        values = [outscale_nat_service.outscale_nat_service.state]
    }
    filter {
        name   = "tags"
        values = ["Key=value-tags"]
   }
}
