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
    policy_orn = outscale_policy.policy-1.orn
  }
}

resource "outscale_user_group" "group-1" {
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
  policy {
    policy_orn = outscale_policy.policy-1.orn
  }

  policy {
    policy_orn = outscale_policy.policy-2.orn
  }
  depends_on = [outscale_user.user-1, outscale_user.user-2]
}

resource "outscale_user_group" "group-2" {
  user_group_name = "test-usergroup-${random_string.suffix[1].result}"
  path            = "/terraform3/"
  user {
    user_name = outscale_user.user-1.user_name
    path      = "/terraform/"
  }
  user {
    user_name = outscale_user.user-2.user_name
    path      = "/terraform2/"
  }
  policy {
    policy_orn = outscale_policy.policy-2.orn
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
  path        = "/terraform2/"
}
