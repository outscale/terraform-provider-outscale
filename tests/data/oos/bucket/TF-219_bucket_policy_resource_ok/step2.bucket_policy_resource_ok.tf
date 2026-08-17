resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_policy" "policy" {
  bucket = outscale_oos_bucket.bucket.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = "*"
      Action    = "s3:PutObject"
      Resource  = "arn:aws:s3:::${outscale_oos_bucket.bucket.id}/*"
    }]
  })
}
