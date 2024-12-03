resource "outscale_policy" "policy_userGroup01" {
  description = "Example of description"
  document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
  path        = "/okht/"
  policy_name = "group-policy"
}

resource "outscale_policy" "policy_userTest" {
  description = "Example of description"
  document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
  path        = "/"
  policy_name = "user-policy"
}

resource "outscale_user" "userTest" {
  user_name = "group_user"
  path = "/IntegGroup/"
  policy {
    policy_orn = outscale_policy.policy_userTest.orn
  }
}

resource "outscale_user_group" "dataUserGroupInteg" {
  user_group_name = "testDataugInteg"
  path            = "/TestdataUG/"
  policy {
    policy_orn = outscale_policy.policy_userGroup01.orn
  }
}

data "outscale_user_groups" "testgrpData" {
  filter {
    name   = "user_group_ids"
    values = [outscale_user_group.dataUserGroupInteg.user_group_id]
  }
  filter {
    name   = "path_prefix"
    values = ["/TestdataUG/"]
  }
}

