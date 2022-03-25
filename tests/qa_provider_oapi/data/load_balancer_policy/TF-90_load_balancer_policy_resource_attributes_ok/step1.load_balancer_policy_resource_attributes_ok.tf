resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name        = "lbu-TF-90"
   subregion_names           = ["${var.region}a"]
   listeners {
      backend_port           = 80
      backend_protocol       = "TCP"
      load_balancer_protocol = "TCP"
      load_balancer_port     = 80
    }
   tags {
      key   = "name"
      value = "public_lbu1"
    }
}

resource "outscale_load_balancer_policy" "policy-1" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name        = "policy-lbu-terraform-1" 
    policy_type        = "load_balancer"
depends_on = [outscale_load_balancer.public_lbu1]
}

resource "outscale_load_balancer_policy" "policy-2" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name        = "policy-lbu-terraform-2"
    policy_type        = "app" 
    cookie_name        = "Cookie-1"
depends_on = [outscale_load_balancer_policy.policy-1]
}

resource "outscale_load_balancer_policy" "policy-3" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name        = "policy-lbu-terraform-3"
    policy_type        = "load_balancer"
depends_on = [outscale_load_balancer_policy.policy-2]
}

resource "outscale_load_balancer_policy" "policy-4" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name        = "policy-lbu-terraform-4"
    policy_type        = "app" 
    cookie_name        = "Cookie-2"
depends_on = [outscale_load_balancer_policy.policy-3]
}
