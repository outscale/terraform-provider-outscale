resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_lifecycle" "lifecycle" {
  bucket     = outscale_oos_bucket.bucket.id

  rule {
    id     = "expiration"
    status = "Enabled"

    expiration = {
      date = "2050-01-01T00:00:00Z"
    }

    filter = {
      prefix = "docs/"
    }
  }

  rule {
    id     = "prefix-tags"
    status = "Disabled"

    expiration = {
      days = 30
    }

    noncurrent_version_expiration = {
      noncurrent_days = 45
    }

    filter = {
      and = {
        prefix = "dir/"
        tags = {
          category = "test"
          owner    = "user"
        }
      }
    }
  }

  rule {
    id     = "abort-uploads"
    status = "Enabled"

    abort_incomplete_multipart_upload = {
      days_after_initiation = 7
    }

    filter = {
      prefix = "uploads/"
    }
  }

  rule {
    status = "Enabled"

    expiration = {
      expired_object_delete_marker = true
    }

    filter = {}
  }
}
