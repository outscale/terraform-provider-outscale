resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_encryption" "encryption" {
  bucket = outscale_oos_bucket.bucket.id
  encryption_type = "AES256"
}
