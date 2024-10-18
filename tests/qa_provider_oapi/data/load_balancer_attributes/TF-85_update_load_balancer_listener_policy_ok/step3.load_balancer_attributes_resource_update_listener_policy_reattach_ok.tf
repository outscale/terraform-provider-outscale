resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name ="lbu-TF-85-1"
   subregion_names= ["${var.region}a"]
   listeners {
      backend_port = 80
      backend_protocol= "HTTP"
      load_balancer_protocol= "HTTP"
      load_balancer_port = 80
      }
   listeners {
      backend_port = 1024
      backend_protocol= "TCP"
      load_balancer_protocol= "TCP"
      load_balancer_port = 1024
      }
   tags {
      key = "name"
      value = "public_lbu1"
      }
   tags {
      key = "test"
      value = "tags"
     }
}

resource "outscale_load_balancer_policy" "policy-1" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name        = "policy-lbu-terraform-TF85"
    policy_type        = "load_balancer"
}


resource "outscale_load_balancer_attributes" "attributes-policy" {
   load_balancer_name = outscale_load_balancer.public_lbu1.id
   load_balancer_port = 80
   policy_names       = ["policy-lbu-terraform-TF85" ]
depends_on =[outscale_load_balancer_policy.policy-1]
}

