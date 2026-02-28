resource "outscale_user" "user-1" {
  user_name  = "test-user-${random_string.suffix[2].result}"
  user_email = "test-TF11@test2.fr"
  path       = "/terraform_update/"
}


resource "outscale_user" "user-2" {
  user_name = "test-user-${random_string.suffix[1].result}"
}

resource "outscale_access_key" "access_key_eim01" {
  user_name       = outscale_user.user-2.user_name
  state           = "INACTIVE"
  expiration_date = "2128-01-01"
  depends_on      = [outscale_user.user-2]
}
