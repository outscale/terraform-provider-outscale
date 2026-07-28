resource "outscale_oos_bucket" "bucket" {
  name         = "test-oos-bucket-${random_string.suffix[0].result}"
  acl          = "public-read"
  force_delete = true
}
