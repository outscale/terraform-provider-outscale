resource "outscale_user" "user-policy"  {
 user_name = "User-TF-linkPolicy"
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
  description = "test-terraform-2"
  document = file("data/policies_files/policy12.json")
  path = "/terraform2/"
}


data "outscale_policies_linked_to_user" "linked_policy01" {

   user_name= outscale_user.user-policy.user_name

depends_on=[outscale_user.user-policy]

}
