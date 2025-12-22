resource "outscale_keypair" "my_keypair" {
 keypair_name = "test-keypair-${random_string.suffix[0].result}"
}

resource "outscale_security_group" "public_sg" {
    description         = "test lbu vm health"
    security_group_name = "test-sg-${random_string.suffix[0].result}"
}

resource "outscale_security_group_rule" "outscale_security_group_rule" {
    flow              = "Inbound"
    security_group_id = outscale_security_group.public_sg.id
    from_port_range   =  8080
    to_port_range     =  8080
    ip_protocol       =  "tcp"
    ip_range          =  "0.0.0.0/0"
}

resource "outscale_vm" "outscale_vm-1" {
   image_id           = var.image_id
   vm_type            = var.vm_type
   keypair_name       = outscale_keypair.my_keypair.keypair_name
   security_group_ids = [outscale_security_group.public_sg.id]
   placement_subregion_name = "${var.region}a"
   user_data = base64encode(<<EOF
      #!/bin/bash
      pushd /home
      nohup python -m SimpleHTTPServer 8080
      EOF
    )
}

resource "outscale_load_balancer" "public_lbu2" {
    load_balancer_name = "test-lb-${random_string.suffix[0].result}"
    subregion_names    = ["${var.region}a"]
    listeners {
       backend_port          = 8080
       backend_protocol      = "HTTP"
       load_balancer_protocol= "HTTP"
       load_balancer_port    = 8080
     }
    tags {
       key   = "name"
       value = "public_lbu2"
    }
}

resource "outscale_load_balancer_attributes" "attributes-1" {
   load_balancer_name      = outscale_load_balancer.public_lbu2.id
    health_check {
        healthy_threshold   = 10
        check_interval      = 30
        port                = 8080
        protocol            = "HTTP"
        timeout             = 5
        unhealthy_threshold = 5
    }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms" {
    load_balancer_name = outscale_load_balancer.public_lbu2.load_balancer_name
    backend_vm_ids     = [outscale_vm.outscale_vm-1.vm_id]
}

data "outscale_load_balancer_vm_health""outscale_load_balancer_vm_health" {
     load_balancer_name  = outscale_load_balancer.public_lbu2.load_balancer_name
     backend_vm_ids      = [outscale_vm.outscale_vm-1.vm_id]
depends_on = [outscale_load_balancer_attributes.attributes-1,outscale_load_balancer_vms.outscale_load_balancer_vms]
}
