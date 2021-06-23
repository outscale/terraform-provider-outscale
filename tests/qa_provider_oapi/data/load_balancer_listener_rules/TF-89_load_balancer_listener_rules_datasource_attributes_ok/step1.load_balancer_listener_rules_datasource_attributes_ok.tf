resource "outscale_keypair" "my_keypair" {
 keypair_name = "KP-TF89"
}
resource "outscale_vm" "public_vm_1" {
   image_id                 = var.image_id
   vm_type                  = var.vm_type
   keypair_name             = outscale_keypair.my_keypair.keypair_name
}

resource "outscale_load_balancer" "public_lbu2" {
   load_balancer_name ="lbu-TF-89"
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
      listener_rule_name      = "listener-rule-1"
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
      listener_rule_name      = "listener-rule-2"
      host_name_pattern       = "*.abc.-.abc.*.com"
      priority                 = 1
    }
   vm_ids = [outscale_vm.public_vm_1.vm_id ]
}

data "outscale_load_balancer_listener_rules" "outscale_load_balancer_listener_rules" {
 filter {
        name = "listener_rule_names"
        values = ["listener-rule-2","listener-rule-1"]
      }
 depends_on =[outscale_load_balancer_listener_rule.rule-2,outscale_load_balancer_listener_rule.rule-1]
}
