resource "outscale_user" "user-1"  {
 user_name = "User-TF-11"
 user_email = "test-TF11@test2.fr"
 path            = "/terraform/"
}

resource "outscale_user" "user-2"  {
 user_name = "User-TF-12"
 user_email = "test-TF12@test2.fr"
 path            = "/terraform2/"
}

resource "outscale_user_group" "group-1" {
 user_group_name = "Group-TF-test-1"
 path            = "/terraform/"
 user {
    user_name = outscale_user.user-1.user_name
    path            = "/terraform/"
 }
  user {
    user_name = outscale_user.user-2.user_name
    path            = "/terraform2/"
 }
depends_on=[outscale_user.user-1,outscale_user.user-2]
}

resource "outscale_user_group" "group-2" {
 user_group_name = "Group-TF-test-2"
 path            = "/terraform3/"
}


data "outscale_user_groups" "usegroups01" {
   filter {
        name   = "user_group_ids"
        values = [outscale_user_group.group-1.user_group_id,outscale_user_group.group-2.user_group_id]
    }
 filter {
        name   = "path_prefix"
        values = [outscale_user_group.group-1.path]
    }
}

data "outscale_user_groups" "usegroups02" {
   filter {
        name   = "user_group_ids"
        values = [outscale_user_group.group-1.user_group_id,outscale_user_group.group-2.user_group_id]
    }
}
