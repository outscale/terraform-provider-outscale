resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF87"
}

resource "outscale_security_group" "my_sgLbl" {
   description = "test sg-group-lbu"
   security_group_name = "SG-inteLbl"
}

resource "outscale_vm" "public_vm_1" {
   image_id           = var.image_id
   vm_type            = var.vm_type
   keypair_name       = outscale_keypair.my_keypair.keypair_name
   security_group_ids = [outscale_security_group.my_sgLbl.security_group_id]
}

resource "outscale_load_balancer" "public_lbu2" {
   load_balancer_name ="lbu-TF-87-${var.suffixe_lbu_name}"
   subregion_names= ["${var.region}a"]
   listeners {
      backend_port = 80
      backend_protocol= "TCP"
      load_balancer_protocol= "TCP"
      load_balancer_port = 80
     }
   tags {
      key = "name"
      value = "listener-rule"
     }
}

resource "outscale_load_balancer_vms" "outscale_load_balancer_vms" {
    load_balancer_name = outscale_load_balancer.public_lbu2.id
    backend_vm_ids = [outscale_vm.public_vm_1.vm_id]
 }


resource "outscale_load_balancer_listener_rule" "rule-1" {
    listener {
       load_balancer_name = outscale_load_balancer.public_lbu2.id
       load_balancer_port = 80
     }

    listener_rule {
      action                  = "forward"
      listener_rule_name      = "listener-rule-1-${var.suffixe_lbu_name}"
      path_pattern             = "*.abc.*.abc.*.com"
      priority                 = 10
    }
   vm_ids = [outscale_vm.public_vm_1.vm_id ]
}

resource "outscale_load_balancer_listener_rule" "rule-2" {
    listener  {
       load_balancer_name = outscale_load_balancer.public_lbu2.id
       load_balancer_port = 80
     }

    listener_rule {
      action                  = "forward"
      listener_rule_name      = "listener-rule-2-${var.suffixe_lbu_name}"
      host_name_pattern       = "*.abc.-.abc.*.com"
      priority                 = 1
    }
   vm_ids = [outscale_vm.public_vm_1.vm_id ]
}

data "outscale_load_balancer_listener_rule" "outscale_load_balancer_listener_rule" {
   filter {
        name = "listener_rule_names"
        values = ["listener-rule-2-${var.suffixe_lbu_name}"]
      }
 depends_on =[outscale_load_balancer_listener_rule.rule-2]
}

