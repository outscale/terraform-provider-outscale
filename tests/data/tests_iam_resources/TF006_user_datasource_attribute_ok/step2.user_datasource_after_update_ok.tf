resource "outscale_user" "user1" {
  user_name  = "test-user-${random_string.suffix[0].result}"
  user_email = "test@test65.fr"
  path       = "/test/terraform/"
}

resource "outscale_access_key" "access_key_01" {
  user_name       = outscale_user.user1.user_name
  expiration_date = "2038-01-01"
  state           = "INACTIVE"
  depends_on      = [outscale_user.user1]
}


data "outscale_user" "user01" {
  filter {
    name   = "user_ids"
    values = [outscale_user.user1.user_id]
  }
}

data "outscale_access_key" "access_key_user01" {
  user_name = outscale_user.user1.user_name
  filter {
    name   = "access_key_ids"
    values = [outscale_access_key.access_key_01.access_key_id]
  }

  filter {
    name   = "states"
    values = ["INACTIVE"]
  }
  depends_on = [outscale_user.user1]
}


data "outscale_access_keys" "access_keys_user01" {
  user_name = outscale_user.user1.user_name
  filter {
    name   = "access_key_ids"
    values = [outscale_access_key.access_key_01.access_key_id]
  }
  depends_on = [outscale_user.user1]
}

