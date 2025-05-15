## Create a Load Balancer###

resource "outscale_load_balancer" "load_balancer01" {
   load_balancer_name = "terraform-lb-TF186-${var.suffixe_lbu_name}"
    subregion_names    = ["${var.region}a"]
    listeners {
        backend_port           = 8080
        backend_protocol       = "HTTP"
        load_balancer_protocol = "HTTP"
        load_balancer_port     = 8080
     }
}


## Create the datasource###

data "outscale_quotas" "lbu-global" {
    filter {
        name   = "quota_types"
        values = ["global"]
      }
  }
