resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name = "lbu-TF-91"
   subregion_names    = ["${var.region}a"]
   listeners {
     backend_port           = 80
     backend_protocol       = "HTTP"
     load_balancer_protocol = "HTTP"
     load_balancer_port     = 80
    }
   tags {
     key = "name"
     value = "public_lbu1"
   }
   tags { 
     key = "Platfotm"
     value = "terraform"
  }
  tags {
    key = "User"
   value = "mzi"
  }
}

resource "outscale_load_balancer" "public_lbu2" {
   load_balancer_name ="lbu-TF-91-2"
   subregion_names= ["${var.region}a"]
   listeners {
      backend_port = 80
      backend_protocol= "HTTP"
      load_balancer_protocol= "HTTP"
      load_balancer_port = 80
    }
   tags {
      key = "name"
      value = "public_lbu2"
   }
   tags {
      key = "Platfotm"
      value = "terraform"
   }
}

data "outscale_load_balancer_tags" "test-tags" {
   load_balancer_names = [outscale_load_balancer.public_lbu1.id,outscale_load_balancer.public_lbu2.id]
}
