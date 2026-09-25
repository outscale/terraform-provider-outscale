resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Enabled"
}

resource "outscale_oos_object" "object-null-version" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "test"
  content = "test"

  depends_on = [outscale_oos_bucket_versioning.versioning]
}

resource "outscale_oos_object" "object-versioned" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "test"
  content = "test"

  depends_on = [outscale_oos_bucket_versioning.versioning]
}
