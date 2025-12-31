resource "outscale_user_group" "group-1" {
   user_group_name = "test-usergroup-${random_string.suffix[1].result}"
   path            = "/"
 }
