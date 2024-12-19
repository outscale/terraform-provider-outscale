resource "outscale_user" "user-1"  {
   user_name = "User-terraform-11"
   user_email = "test-TF11@test2.fr"
   path            = "/terraform_update/"
 }


resource "outscale_user" "user-2"  {
   user_name = "User-terraform-2"
 }

resource "outscale_access_key" "access_key_eim01" {
    user_name = outscale_user.user-2.user_name
    state           = "ACTIVE"
    expiration_date = "2028-01-01"
depends_on=[outscale_user.user-2]
}

