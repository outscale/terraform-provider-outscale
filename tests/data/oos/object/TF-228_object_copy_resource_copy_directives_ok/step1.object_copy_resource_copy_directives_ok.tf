resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-${random_string.suffix[0].result}"
}

resource "outscale_oos_object" "source" {
  bucket              = outscale_oos_bucket.bucket.id
  key                 = "source"
  content             = "content"
  cache_control       = "no-store"
  content_disposition = "attachment"
  content_encoding    = "gzip"
  content_language    = "fr"
  content_type        = "application/json"
  expires             = "2035-05-06T07:08:09Z"
  metadata = {
    key = "value"
  }
  tags = {
    key = "value"
  }
}

resource "outscale_oos_object_copy" "copy" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination"
  copy_source = {
    path = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
  }
  metadata_directive = "COPY"
  tagging_directive  = "COPY"
}
