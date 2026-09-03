resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Enabled"
}

resource "outscale_oos_object" "object" {
  bucket     = outscale_oos_bucket.bucket.id
  key        = "versioned-object.txt"
  content    = "test versioned object content"
  depends_on = [outscale_oos_bucket_versioning.versioning]
  tags = {
    versioned = "true"
  }
}
