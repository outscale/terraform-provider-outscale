resource "outscale_policy" "policy-1"  {
  policy_name = "terraform-policy-11"
  description = "test-terraform-11"
  document = file("data/policies_files/policy11.json")
  path = "/policy1/"
}
resource "outscale_policy_version" "policy11-version-02" {
  policy_orn = outscale_policy.policy-1.orn
  document = file("data/policies_files/policy12.json")
  set_as_default = true
}

resource "outscale_user" "user-1"  {
 user_name = "User-TF-1"
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }
}

resource "outscale_user_group" "group-1" {
 user_group_name = "Group-TF-test-1"
 path            = "/terraform/"
policy {
  policy_orn = outscale_policy.policy-1.orn
 }
}

resource "outscale_policy" "policy-12"  {
  policy_name = "terraform-policy-12"
  description = "test-terraform-12"
  document = file("data/policies_files/policy12.json")
  path = "/policy12/"
}


data "outscale_policies" "Mypolicies01" {
 filter {
     name = "path_prefix"
     values = [outscale_policy.policy-1.path,outscale_policy.policy-12.path]
    }
}
