resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_lifecycle" "lifecycle" {
  bucket     = outscale_oos_bucket.bucket.id

  rule {
    id     = "expiration"
    status = "Disabled"

    expiration = {
      date = "2051-02-02T00:00:00Z"
    }

    filter = {
      prefix = "docs/"
    }
  }

  rule {
    id     = "prefix-tags"
    status = "Enabled"

    expiration = {
      days = 60
    }

    noncurrent_version_expiration = {
      noncurrent_days = 90
    }

    filter = {
      and = {
        prefix = "dir/"
        tags = {
          category = "test"
          team     = "user"
        }
      }
    }
  }

  rule {
    id     = "abort-uploads"
    status = "Disabled"

    abort_incomplete_multipart_upload = {
      days_after_initiation = 14
    }

    filter = {
      prefix = "uploads/"
    }
  }

  rule {
    status = "Disabled"

    expiration = {
      expired_object_delete_marker = false
    }

    filter = {}
  }
}
