resource "outscale_oos_bucket" "bucket" {
  name         = "test-oos-bucket-${random_string.suffix[0].result}"
  acl          = "private"
  force_delete = false
}
