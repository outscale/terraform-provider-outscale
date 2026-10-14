resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-${random_string.suffix[0].result}"
}

resource "outscale_oos_object" "source" {
  bucket       = outscale_oos_bucket.bucket.id
  key          = "source"
  content      = "content"
  content_type = "text/plain"
  metadata = {
    key = "source"
  }
  tags = {
    key = "source"
  }
}

resource "outscale_oos_object_copy" "copy" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination"
  acl    = "public-read"
  copy_source = {
    path = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
  }
  metadata_directive  = "REPLACE"
  tagging_directive   = "REPLACE"
  cache_control       = "no-transform"
  content_disposition = "attachment"
  content_encoding    = "br"
  content_language    = "de"
  content_type        = "application/json"
  expires             = "2036-06-07T08:09:10Z"
  metadata = {
    key = "destination"
  }
  encryption_type = "AES256"
  tags = {
    key = "value2"
  }
}
