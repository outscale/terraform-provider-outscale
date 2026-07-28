resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-bucket-${random_string.suffix[0].result}"
  object_lock = {
    enabled = true
    default_retention = {
      days = 2
    }
  }
}
