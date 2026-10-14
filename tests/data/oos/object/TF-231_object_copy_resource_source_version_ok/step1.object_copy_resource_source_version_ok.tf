resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-${random_string.suffix[0].result}"
}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Enabled"
}

resource "outscale_oos_object" "source1" {
  bucket     = outscale_oos_bucket.bucket.id
  key        = "source"
  content    = "v1"
  depends_on = [outscale_oos_bucket_versioning.versioning]
}
