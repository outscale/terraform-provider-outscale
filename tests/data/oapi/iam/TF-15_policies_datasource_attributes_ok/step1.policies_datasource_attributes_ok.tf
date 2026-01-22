resource "outscale_policy" "policy-1" {
  policy_name = "test-policy-${random_string.suffix[0].result}"
  description = "test-terraform-11"
  document    = file("policies/policy11.json")
  path        = "/policy1/"
}
resource "outscale_policy_version" "policy11-version-02" {
  policy_orn     = outscale_policy.policy-1.orn
  document       = file("policies/policy12.json")
  set_as_default = true
}

resource "outscale_user" "user-1" {
  user_name = "test-user-${random_string.suffix[0].result}"
  policy {
    policy_orn = outscale_policy.policy-1.orn
  }
}

resource "outscale_user_group" "group-1" {
  user_group_name = "test-usergroup-${random_string.suffix[0].result}"
  path            = "/terraform/"
  policy {
    policy_orn = outscale_policy.policy-1.orn
  }
}

resource "outscale_policy" "policy-12" {
  policy_name = "test-policy-${random_string.suffix[1].result}"
  description = "test-terraform-12"
  document    = file("policies/policy12.json")
  path        = "/policy12/"
}


data "outscale_policies" "Mypolicies01" {
  filter {
    name   = "path_prefix"
    values = [outscale_policy.policy-1.path, outscale_policy.policy-12.path]
  }
}
