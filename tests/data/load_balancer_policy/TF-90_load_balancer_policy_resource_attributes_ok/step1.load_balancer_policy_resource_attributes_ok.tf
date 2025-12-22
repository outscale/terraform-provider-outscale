resource "outscale_load_balancer" "public_lbu1" {
   load_balancer_name = "test-lb-${random_string.suffix[0].result}"
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
    policy_name = "test-policy-${random_string.suffix[0].result}"
    policy_type        = "load_balancer"
depends_on = [outscale_load_balancer.public_lbu1]
}

resource "outscale_load_balancer_policy" "policy-2" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name = "test-policy-${random_string.suffix[1].result}"
    policy_type        = "app"
    cookie_name = "test-cookie-${random_string.suffix[0].result}"
depends_on = [outscale_load_balancer_policy.policy-1]
}

resource "outscale_load_balancer_policy" "policy-3" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name = "test-policy-${random_string.suffix[2].result}"
    policy_type        = "load_balancer"
depends_on = [outscale_load_balancer_policy.policy-2]
}

resource "outscale_load_balancer_policy" "policy-4" {
    load_balancer_name = outscale_load_balancer.public_lbu1.load_balancer_name
    policy_name = "test-policy-${random_string.suffix[3].result}"
    policy_type        = "app"
    cookie_name = "test-cookie-${random_string.suffix[1].result}"
depends_on = [outscale_load_balancer_policy.policy-3]
}
