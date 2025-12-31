resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}
resource "outscale_load_balancer" "public_lbu1" {
  load_balancer_name = "test-lb-${random_string.suffix[0].result}"
  subregion_names           = ["${var.region}a"]
  listeners {
     backend_port           = 80
     backend_protocol       = "TCP"
     load_balancer_protocol = "TCP"
     load_balancer_port     = 80
    }
  listeners {
     backend_port           = 8080
     backend_protocol       = "HTTP"
     load_balancer_protocol = "HTTP"
     load_balancer_port     = 8080
    }
  tags {
     key                    = "name"
     value                  = "public_lbu1"
  }
}

resource "outscale_security_group" "public_sglb" {
    description             = "test lbu vm health"
    security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_vm" "outscale_vms_lbu" {
   count                    = 2
   image_id                 = var.image_id
   vm_type                  = var.vm_type
   keypair_name             = outscale_keypair.my_keypair.keypair_name
   security_group_ids       = [outscale_security_group.public_sglb.id]
   user_data                = base64encode(<<EOF
     #!/bin/bash
    pushd /home
    nohup python -m SimpleHTTPServer 8080
    EOF
    )
  tags {
     key                    = "Name"
     value                  = "Backend-Vms-mzi"
  }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms" {
    load_balancer_name      = outscale_load_balancer.public_lbu1.load_balancer_name
    backend_vm_ids          = [outscale_vm.outscale_vms_lbu[0].vm_id]
}
