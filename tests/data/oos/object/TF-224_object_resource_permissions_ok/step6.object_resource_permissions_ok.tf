resource "outscale_oos_bucket" "bucket" {
  name = "test-oos-object-${random_string.suffix[0].result}"
}

resource "outscale_oos_bucket" "owner" {}

resource "outscale_oos_object" "object" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "object.txt"
  content = "test object content"
  grant = {
    full_control = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }
    read = {
      ids             = [outscale_oos_bucket.owner.permissions.owner.id]
      email_addresses = ["customer-tooling@outscale.com"]
    }
    read_acp = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }
    write = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }
    write_acp = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }
  }
}
