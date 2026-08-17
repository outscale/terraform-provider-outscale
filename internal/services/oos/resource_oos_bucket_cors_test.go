package oos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_BucketCors(t *testing.T) {
	resourceName := "outscale_oos_bucket_cors.cors"
	bucketResourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketCorsConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "id", bucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"max_age_seconds": "3000",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "rule.*.allowed_headers.*", "Authorization"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rule.*.allowed_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rule.*.allowed_origins.*", "https://www.example.com"),
				),
			},
			{
				Config: testAccOOSBucketCorsConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"max_age_seconds": "600",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "rule.*", map[string]string{
						"max_age_seconds": "1200",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "rule.*.allowed_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(resourceName, "rule.*.expose_headers.*", "Content-Type"),
				),
			},
		},
	})
}

var testAccOOSBucketCorsConfigBasic = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_cors" "cors" {
  bucket = outscale_oos_bucket.bucket.id

  rule {
    allowed_headers = ["Authorization"]
    allowed_methods = ["GET"]
    allowed_origins = ["https://www.example.com"]
    max_age_seconds = 3000
  }
}
`

var testAccOOSBucketCorsConfigUpdated = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_cors" "cors" {
  bucket = outscale_oos_bucket.bucket.id

  rule {
    allowed_headers = ["Authorization", "Content-Type"]
    allowed_methods = ["GET", "PUT"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["Content-Type"]
    max_age_seconds = 600
  }

  rule {
    allowed_headers = ["*"]
    allowed_methods = ["HEAD"]
    allowed_origins = ["https://www.foobar.example"]
    max_age_seconds = 1200
  }
}
`
