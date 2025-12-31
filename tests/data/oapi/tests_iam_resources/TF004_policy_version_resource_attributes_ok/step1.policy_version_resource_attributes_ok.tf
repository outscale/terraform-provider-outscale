resource "outscale_policy" "policy-1" {
  policy_name = "test-policy-${random_string.suffix[0].result}"
  description = "test-version"
  document    = file("policies/policyV1.json")
  path        = "/"
}

resource "outscale_policy_version" "policy-1_version-02" {
  policy_orn = outscale_policy.policy-1.orn
  document   = file("policies/policyV2.json")
}

resource "outscale_policy_version" "version-03" {
  policy_orn     = outscale_policy.policy-1.orn
  document       = file("policies/policyV3.json")
  set_as_default = true
  depends_on     = [outscale_policy_version.policy-1_version-02]
}
