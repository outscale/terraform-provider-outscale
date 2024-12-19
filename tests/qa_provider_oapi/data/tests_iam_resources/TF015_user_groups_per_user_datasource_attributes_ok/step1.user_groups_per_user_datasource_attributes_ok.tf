resource "outscale_user" "user-1"  {
 user_name = "User-TF-group1"
 path      = "/terraform/"
}



resource "outscale_user_group" "group-1" {
 user_group_name = "Group-TF-test-1"
 path            = "/terraform/"
  user {
    user_name = outscale_user.user-1.user_name
    path            = "/terraform/"
 }
depends_on=[outscale_user.user-1]
}



resource "outscale_user_group" "group-2" {
 user_group_name = "Group-TF-test-2"
 path            = "/terraform3/"
 user {
    user_name = outscale_user.user-1.user_name
    path            = "/terraform/"
 }
depends_on=[outscale_user.user-1]
}

data "outscale_user_groups_per_user" "usegroups_per_user01" {
  user_name = outscale_user.user-1.user_name
  user_path = outscale_user.user-1.path
}
