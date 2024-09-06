resource "outscale_user_group" "dataUserGroupInteg" {
  user_group_name = "testDataugInteg"
  path            = "/TestdataUG/"
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
