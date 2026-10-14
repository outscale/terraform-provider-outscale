resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-${random_string.suffix[0].result}"
}

resource "outscale_oos_object" "source" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "source"
  content = "content"
}

resource "outscale_oos_object_copy" "copy1" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination1"
  copy_source = {
    path     = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
    if_match = outscale_oos_object.source.etag
  }
}

resource "outscale_oos_object_copy" "copy2" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination2"
  copy_source = {
    path          = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
    if_none_match = "\"00000000000000000000000000000000\""
  }
}

resource "outscale_oos_object_copy" "copy3" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination3"
  copy_source = {
    path              = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
    if_modified_since = "2000-01-01T00:00:00Z"
  }
}

resource "outscale_oos_object_copy" "copy4" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination4"
  copy_source = {
    path                = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
    if_unmodified_since = "2099-01-01T00:00:00Z"
  }
}
