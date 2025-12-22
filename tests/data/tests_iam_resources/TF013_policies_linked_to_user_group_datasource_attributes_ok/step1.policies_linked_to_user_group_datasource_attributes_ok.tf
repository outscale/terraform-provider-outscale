resource "outscale_user_group" "group-1" {
 user_group_name = "Group-TF-policy-1"
 path            = "/terraform/"
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }

 policy {
  policy_orn = outscale_policy.policy-2.orn
 }
}

resource "outscale_policy" "policy-1"  {
  policy_name = "terraform-policy-1"
  description = "test-terraform"
  document = file("data/policies_files/policy11.json")
  path = "/"
}

resource "outscale_policy" "policy-2"  {
  policy_name = "terraform-policy-2"
  description = "test-terraform"
  document = file("data/policies_files/policy12.json")
  path = "/terraform2/"
}

data "outscale_policies_linked_to_user_group" "managed_policies_linked_to_user_group" {
   user_group_name= outscale_user_group.group-1.user_group_name
depends_on=[outscale_user_group.group-1]
}
