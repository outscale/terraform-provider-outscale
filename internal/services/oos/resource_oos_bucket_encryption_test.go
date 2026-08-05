package oos_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_Bucket_Encryption(t *testing.T) {
	resourceName := "outscale_oos_bucket_encryption.encryption"
	bucketResourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketEncryptionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "encryption_type", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "bucket_key_enabled", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "id", bucketResourceName, "id"),
				),
			},
		},
	})
}

var testAccOOSBucketEncryptionConfig = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_encryption" "encryption" {
  bucket = outscale_oos_bucket.bucket.id
  encryption_type = "AES256"
}
`
