resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/24"
    tags {
     key = "name"
     value = "Net"
    }
}

resource "outscale_route_table" "outscale_route_table" {
    net_id = outscale_net.outscale_net.net_id
    #depends_on = ["outscale_internet_service_link.outscale_internet_service_link"]
     tags {
     key = "name"
     value = "Route-Table"
    }
}

resource "outscale_internet_service" "outscale_internet_service" {
 tags {
     key = "name"
     value = "InternetService"
    }
}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
   # internet_service_id = outscale_internet_service.outscale_internet_service.id    # TEST purpose only
    net_id = outscale_net.outscale_net.net_id
}


resource "outscale_route" "outscale_route" {
    gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
    #gateway_id           = outscale_internet_service.outscale_internet_service.id    # TEST purposes only
    destination_ip_range = "10.0.0.0/16"
    route_table_id       = outscale_route_table.outscale_route_table.route_table_id
    depends_on = [outscale_internet_service_link.outscale_internet_service_link]
}
