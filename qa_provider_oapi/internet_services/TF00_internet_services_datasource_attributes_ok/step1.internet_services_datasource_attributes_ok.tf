resource "outscale_internet_service" "outscale_internet_service" {
}

resource "outscale_internet_service" "outscale_internet_service2" {
}

data "outscale_internet_services" "outscale_internet_services" {
    filter {
        name   = "internet_service_ids"
        values = [outscale_internet_service.outscale_internet_service.internet_service_id, outscale_internet_service.outscale_internet_service2.internet_service_id]
    }
}
