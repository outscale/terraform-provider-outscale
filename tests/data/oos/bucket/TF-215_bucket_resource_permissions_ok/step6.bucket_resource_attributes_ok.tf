resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-bucket-${random_string.suffix[0].result}"
  grant = {
    read = {
      ids             = [outscale_oos_bucket.owner.permissions.owner.id]
      email_addresses = ["customer-tooling@outscale.com"]
    }
    full_control = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }
  }
  force_delete = false
}

resource "outscale_oos_bucket" "owner" {}
