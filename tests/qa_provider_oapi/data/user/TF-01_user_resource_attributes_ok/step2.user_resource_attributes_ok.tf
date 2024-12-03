resource "outscale_policy" "policy_user01" {
  description = "Example of description"
  document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
  path        = "/okht/"
  policy_name = "okht-user-policy"
}

resource "outscale_user" "userInteg" {
  user_name = "test_integ_update"
  path = "/Integ/"
  policy {
    policy_orn = outscale_policy.policy_user01.orn
  }
}
