resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_policy" "policy" {
  bucket = outscale_oos_bucket.bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Deny"
      Principal = "*"
      Action    = "s3:GetObject"
      Resource  = "arn:aws:s3:::${outscale_oos_bucket.bucket.id}/*"
    }]
  })
}
