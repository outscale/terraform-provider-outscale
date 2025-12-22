resource "outscale_user" "user-1"  {
 user_name = "User-TF-11"
 user_email = "test-TF11@test2.fr"
 path            = "/terraform/"
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }
}

resource "outscale_user" "user-2"  {
 user_name = "User-TF-12"
 user_email = "test-TF12@test2.fr"
 path            = "/terraform2/"
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }
}

resource "outscale_user_group" "group-1" {
 user_group_name = "Group-TF-test-1"
 path            = "/terraform/"
 user {
    user_name = outscale_user.user-1.user_name
    path            = "/terraform/"
 }
  user {
    user_name = outscale_user.user-2.user_name
    path            = "/terraform2/"
 }
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }

 policy {
  policy_orn = outscale_policy.policy-2.orn
 }
depends_on=[outscale_user.user-1,outscale_user.user-2]
}

resource "outscale_user_group" "group-2" {
 user_group_name = "Group-TF-test-2"
 path            = "/terraform3/"
 user {
    user_name = outscale_user.user-1.user_name
    path            = "/terraform/"
 }
  user {
    user_name = outscale_user.user-2.user_name
    path            = "/terraform2/"
 }
policy {
  policy_orn = outscale_policy.policy-2.orn
 }
depends_on=[outscale_user.user-1,outscale_user.user-2]
}

resource "outscale_policy" "policy-1"  {
  policy_name = "terraform-policy-1"
  description = "test-terraform"
  document = file("data/policies_files/policy.json")
  path = "/"
}

resource "outscale_policy" "policy-2"  {
  policy_name = "terraform-policy-2"
  description = "test-terraform"
  document = file("data/policies_files/policy.json")
  path = "/terraform2/"
}


