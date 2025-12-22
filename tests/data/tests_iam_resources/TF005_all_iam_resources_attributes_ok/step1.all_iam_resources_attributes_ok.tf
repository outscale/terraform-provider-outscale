resource "outscale_user" "user-1" {
  user_name  = "test-user-${random_string.suffix[0].result}"
  user_email = "test-TF11@test2.fr"
  path       = "/terraform/"
  policy {
    policy_orn = outscale_policy.policy-1.orn
  }
}

resource "outscale_user" "user-2" {
  user_name  = "test-user-${random_string.suffix[1].result}"
  user_email = "test-TF12@test2.fr"
  path       = "/terraform2/"
  policy {
    policy_orn         = outscale_policy.policy-2.orn
    default_version_id = "V2"
  }
}


resource "outscale_user_group" "group-No-policy" {
  user_group_name = "test-usergroup-${random_string.suffix[0].result}"
  path            = "/terraform/"
  user {
    user_name = outscale_user.user-1.user_name
    path      = "/terraform/"
  }
  user {
    user_name = outscale_user.user-2.user_name
    path      = "/terraform2/"
  }
  depends_on = [outscale_user.user-1, outscale_user.user-2]
}


resource "outscale_policy" "policy-1" {
  policy_name = "test-policy-${random_string.suffix[0].result}"
  description = "test-terraform"
  document    = file("policies/policy.json")
  path        = "/"
}

resource "outscale_policy" "policy-2" {
  policy_name = "test-policy-${random_string.suffix[1].result}"
  description = "test-terraform"
  document    = file("policies/policy.json")
  path        = "/"
}

resource "outscale_policy_version" "policy-2" {
  policy_orn = outscale_policy.policy-2.orn
  document   = file("policies/policy2.json")
}

#################################################################@

#####Create the users belonging to the same policy group####
resource "outscale_user" "multiple_users" {
  count      = 3
  user_name  = "test-user-${random_string.suffix[2].result}-${count.index}"
  user_email = "test-${count.index}@test2.fr"
}



resource "outscale_user_group" "group-policy-RO" {
  user_group_name = "test-usergroup-${random_string.suffix[1].result}"
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
  depends_on = [outscale_user.multiple_users[0], outscale_user.multiple_users[1], outscale_user.multiple_users[2]]
}


resource "outscale_policy" "policy-RO" {
  policy_name = "test-policy-${random_string.suffix[2].result}"
  description = "test-terraform-ro"
  document    = file("policies/policy_RO.json")
  path        = "/"
}
