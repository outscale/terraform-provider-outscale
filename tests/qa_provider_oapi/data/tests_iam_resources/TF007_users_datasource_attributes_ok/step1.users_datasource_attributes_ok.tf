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

data "outscale_users" "my_users" {
   filter {
        name   = "user_ids"
        values = [outscale_user.user-1.user_id,outscale_user.user-2.user_id]
    }
}

