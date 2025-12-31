resource "outscale_user" "user-1" {
  user_name  = "test-user-${random_string.suffix[0].result}"
  user_email = "test-TF011@test2.fr"
  path       = "/terraform/"
}

resource "outscale_user" "user-2" {
  user_name  = "test-user-${random_string.suffix[1].result}"
  user_email = "test-TF002@test2.fr"
  path       = "/terraform2/"
}



resource "outscale_user_group" "group-1" {
  user_group_name = "test-usergroup-${random_string.suffix[0].result}"
  path            = "/terraform/"
  user {
    user_name = outscale_user.user-1.user_name
    path      = "/terraform/"
  }
  user {
    user_name = outscale_user.user-2.user_name
    path      = "/terraform2/"
  }
  depends_on = [outscale_user.user-1, outscale_user.user-2]
}



data "outscale_user_group" "user_group01" {
  user_group_name = outscale_user_group.group-1.user_group_name
  path            = outscale_user_group.group-1.path
  depends_on      = [outscale_user_group.group-1]
}
