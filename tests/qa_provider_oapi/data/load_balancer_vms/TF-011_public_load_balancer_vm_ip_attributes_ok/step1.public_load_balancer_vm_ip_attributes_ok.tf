resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF011"
}

resource "outscale_security_group" "security_groupTF011" {
    description          = "Terraform"
    security_group_name = "terraform-TF011"
}


resource "outscale_vm" "vm011" {
    count              = 3
    security_group_ids = [outscale_security_group.security_groupTF011.security_group_id]
    image_id           = var.image_id
    vm_type            = var.vm_type
    keypair_name       = outscale_keypair.my_keypair.keypair_name
tags {
     key                    = "name"
     value                  = "vm-TF011"
  }
}


resource "outscale_load_balancer" "public_lbu11" {
  load_balancer_name        = "lbu-TF-011-${var.suffixe_lbu_name}"
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


resource "outscale_load_balancer_vms" "outscale_load_balancer_vms011" {
    load_balancer_name = outscale_load_balancer.public_lbu11.load_balancer_name
    backend_ips     = [outscale_vm.vm011[0].public_ip,outscale_vm.vm011[1].public_ip]
}

data "outscale_load_balancer" "load_balancer_TF011" {
    filter {
        name   = "load_balancer_names"
        values = [outscale_load_balancer.public_lbu11.load_balancer_name]
    }
depends_on =[outscale_load_balancer_vms.outscale_load_balancer_vms011]
}
