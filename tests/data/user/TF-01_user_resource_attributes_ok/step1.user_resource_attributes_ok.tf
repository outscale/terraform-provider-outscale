resource "outscale_user" "userInteg" {
  user_name = "test-user-${random_string.suffix[0].result}"
  path = "/"
}
