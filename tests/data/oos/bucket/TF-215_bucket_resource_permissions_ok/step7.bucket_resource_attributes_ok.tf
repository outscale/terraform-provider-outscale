resource "outscale_oos_bucket" "bucket" {
  name         = "test-oos-bucket-${random_string.suffix[0].result}"
  force_delete = false
}

resource "outscale_oos_bucket" "owner" {}
