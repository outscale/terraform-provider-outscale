package oos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_BucketPolicy(t *testing.T) {
	resourceName := "outscale_oos_bucket_policy.policy"
	bucketResourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "id", bucketResourceName, "id"),
				),
			},
		},
	})
}

var testAccOOSBucketPolicyConfig = `
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
`
