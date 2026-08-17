resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_cors" "cors" {
  bucket = outscale_oos_bucket.bucket.id

  rule {
    allowed_headers = ["xxx", "yyy"]
    allowed_methods = ["PUT", "POST", "DELETE"]
    allowed_origins = ["https://www.example.com", "https://www.foobar.example"]
    expose_headers= ["Content-Type"]
    max_age_seconds = 3001
  }
}
