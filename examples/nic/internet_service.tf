resource "outscale_internet_service" "my_internet_service" {
}

resource "outscale_internet_service_link" "my_internet_service_link" {
  internet_service_id = outscale_internet_service.my_internet_service.internet_service_id
  net_id              = outscale_net.my_net.net_id
}
