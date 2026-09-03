package oos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_Bucket_Versioning(t *testing.T) {
	resourceName := "outscale_oos_bucket_versioning.versioning"
	bucketResourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketVersioningConfig("Enabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "id", bucketResourceName, "id"),
				),
			},
			{
				Config: testAccOOSBucketVersioningConfig("Suspended"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "Suspended"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
				),
			},
			{
				Config: testAccOOSBucketVersioningConfig("Enabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "Enabled"),
				),
			},
		},
	})
}

func testAccOOSBucketVersioningConfig(status string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_bucket_versioning" "versioning" {
  bucket = outscale_oos_bucket.bucket.id
  status = %q
}
`, status)
}
