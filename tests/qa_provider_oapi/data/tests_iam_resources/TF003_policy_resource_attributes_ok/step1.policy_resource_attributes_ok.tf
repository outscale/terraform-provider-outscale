resource "outscale_policy" "policy-1"  {
  policy_name = "terraform-RO-policy"
  description = "test-policy"
  document = file("data/policies_files/policy_RO.json")
  path = "/terraform/"
}

resource "outscale_policy" "policy-2"  {
  policy_name = "terraform-Read-Account_consumption_only"
  document = file("data/policies_files/policy_ReadAccountConsumption.json")
}
