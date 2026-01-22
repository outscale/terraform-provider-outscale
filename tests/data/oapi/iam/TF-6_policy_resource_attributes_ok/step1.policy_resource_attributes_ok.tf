resource "outscale_policy" "policy-1"  {
  policy_name = "test-policy-${random_string.suffix[0].result}"
  description = "test-policy"
  document = file("policies/policy_RO.json")
  path = "/terraform/"
}

resource "outscale_policy" "policy-2"  {
  policy_name = "test-policy-${random_string.suffix[1].result}"
  document = file("policies/policy_ReadAccountConsumption.json")
}
