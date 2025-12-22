resource "outscale_policy" "policy-1" {
  policy_name = "terraform-policy-01"
  description = "test-version"
  document    = file("data/policies_files/policyV1.json")
  path        = "/"
}

resource "outscale_policy_version" "policy-1_version-02" {
  policy_orn = outscale_policy.policy-1.orn
  document   = file("data/policies_files/policyV2.json")
}

resource "outscale_policy_version" "version-03" {
  policy_orn     = outscale_policy.policy-1.orn
  document       = file("data/policies_files/policyV3.json")
  set_as_default = true
  depends_on     = [outscale_policy_version.policy-1_version-02]
}
