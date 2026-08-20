package oos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_BucketLifecycle(t *testing.T) {
	resourceName := "outscale_oos_bucket_lifecycle.lifecycle"
	bucketResourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketLifecycleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "id", bucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				Config: testAccOOSBucketLifecycleConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "4"),
				),
			},
		},
	})
}

var testAccOOSBucketLifecycleConfigBasic = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Enabled"
}

resource "outscale_oos_bucket_lifecycle" "lifecycle" {
  bucket     = outscale_oos_bucket.bucket.id
  depends_on = [outscale_oos_bucket_versioning.versioning]

  rule {
    id     = "expiration"
    status = "Enabled"

    expiration = {
      days = 2
    }

    filter = {}
  }
}
`

var testAccOOSBucketLifecycleConfigUpdated = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = "Enabled"
}

resource "outscale_oos_bucket_lifecycle" "lifecycle" {
  bucket     = outscale_oos_bucket.bucket.id
  depends_on = [outscale_oos_bucket_versioning.versioning]

  rule {
    id     = "expire-on-date"
    status = "Enabled"

    expiration = {
      date = "2050-01-01T00:00:00Z"
    }

    filter = {
      prefix = "dir/"
    }
  }

  rule {
    id     = "expire-noncurrent"
    status = "Disabled"

    noncurrent_version_expiration = {
      noncurrent_days = 30
    }

    filter = {
      and = {
        prefix = ""

        tags = {
          category = "test"
        }
      }
    }
  }

  rule {
    id     = "abort-multipart"
    status = "Enabled"

    abort_incomplete_multipart_upload = {
      days_after_initiation = 7
    }

    filter = {}
  }

  rule {
    status = "Enabled"

    expiration = {
      expired_object_delete_marker = true
    }

    filter = {}
  }
}
`
