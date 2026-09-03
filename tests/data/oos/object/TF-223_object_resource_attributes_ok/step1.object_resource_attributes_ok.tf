resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_object" "object" {
  bucket              = outscale_oos_bucket.bucket.id
  key                 = "directory/object.txt"
  content             = "test object content"
  acl                 = "private"
  cache_control       = "no-cache"
  content_disposition = "inline"
  content_encoding    = "identity"
  content_language    = "en"
  content_type        = "text/plain"
  expires             = "2030-01-01T00:00:00Z"
  metadata = {
    environment = "test"
    owner       = "terraform"
  }
  encryption_type = "AES256"
  tags = {
    environment = "test"
    owner       = "terraform"
  }
}
