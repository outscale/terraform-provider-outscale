
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/18"
}

resource "outscale_internet_service" "outscale_internet_service" {
tags {
      key = "Name"
      value = "Terraform-IGW"
     }
}

resource "outscale_internet_service_link" "outscale_internet_service_link" {
    internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
    net_id              = outscale_net.outscale_net.net_id
}
