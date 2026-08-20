resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_lifecycle" "lifecycle" {
  bucket = outscale_oos_bucket.bucket.id

  rule {
    status = "Disabled"

    expiration = {
      days = 1
    }

    filter = {}
  }
}
