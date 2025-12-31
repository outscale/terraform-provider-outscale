resource "outscale_policy" "policy_user01" {
  description = "Example of description"
  document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
  path        = "/okht/"
  policy_name = "test-policy-${random_string.suffix[0].result}"
}

resource "outscale_user" "userInteg" {
  user_name = "test-user-${random_string.suffix[0].result}"
  path = "/Integ/"
  policy {
    policy_orn = outscale_policy.policy_user01.orn
  }
}
