resource "outscale_user" "user-1"  {
 user_name = "User-TF-RW"
 user_email = "test-TF11@test2.fr"
 path            = "/terraform/"
 policy {
  policy_orn = outscale_policy.policy-1.orn
 }
}

resource "outscale_user" "user-2"  {
 user_name = "User-TF-RW-2"
 user_email = "test-TF12@test2.fr"
 path            = "/terraform2/"
 policy {
  policy_orn = outscale_policy.policy-2.orn
  default_version_id  = "V2"
 }
}


resource "outscale_user_group" "group-No-policy" {
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
  path = "/"
}

resource "outscale_policy_version" "policy-2"  {
  policy_orn = outscale_policy.policy-2.orn
  document = file("data/policies_files/policy2.json")
}

#################################################################@

#####Create the users belonging to the same policy group####
resource "outscale_user" "multiple_users"  {
 count=3
 user_name = "TF-User-${count.index}"
 user_email = "test-${count.index}@test2.fr"
}



resource "outscale_user_group" "group-policy-RO" {
 user_group_name = "Group-TF-RO"
  user {
    user_name = outscale_user.multiple_users[0].user_name
 }
  user {
    user_name = outscale_user.multiple_users[1].user_name
 }
 user {
    user_name = outscale_user.multiple_users[2].user_name
 }
policy {
  policy_orn = outscale_policy.policy-RO.orn
 }
depends_on=[outscale_user.multiple_users[0],outscale_user.multiple_users[1],outscale_user.multiple_users[2]]
}


resource "outscale_policy" "policy-RO"  {
  policy_name = "terraform-policy-RO"
  description = "test-terraform-ro"
  document = file("data/policies_files/policy_RO.json")
  path = "/"
}

