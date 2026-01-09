resource "outscale_internet_service" "outscale_internet_service" {
tags {
     key = "test"
     value = "internet_service"
     }
}

data "outscale_internet_service" "outscale_internet_serviced" {
  filter {
        name   = "tag_keys"
        values = ["test"]
    }
filter {
        name   = "tag_values"
        values = ["internet_service"]
    }
filter {
        name   = "tags"
        values = ["test=internet_service"]
    }
depends_on = [outscale_internet_service.outscale_internet_service]
}
