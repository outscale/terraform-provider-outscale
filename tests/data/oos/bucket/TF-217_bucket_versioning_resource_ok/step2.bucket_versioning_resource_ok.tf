resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Suspended"
}
