resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_bucket" "owner" {}

resource "outscale_oos_object" "object" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "object.txt"
  content = "test object content"
}
