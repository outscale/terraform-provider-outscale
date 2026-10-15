resource "outscale_oos_bucket" "bucket" {
  name         = "test-oos-${random_string.suffix[0].result}"
  force_delete = true
  object_lock = {
    enabled = true
  }
}

resource "outscale_oos_object" "object" {
  bucket       = outscale_oos_bucket.bucket.id
  key          = "object"
  content      = "content"
  force_delete = true
}

resource "outscale_oos_object_retention" "retention" {
  bucket            = outscale_oos_object.object.bucket
  key               = outscale_oos_object.object.key
  version_id        = outscale_oos_object.object.version_id
  retain_until_date = timeadd(timestamp(), "5s")
}

resource "terraform_data" "wait_for_retention" {
  depends_on = [outscale_oos_object_retention.retention]

  provisioner "local-exec" {
    command = "sleep 6"
  }
}
