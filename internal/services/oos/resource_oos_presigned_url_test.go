package oos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_PresignedURL(t *testing.T) {
	bucketName := acctest.RandomWithPrefix("testacc-oos-presigned-url")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccPresignedURLConfig(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("outscale_oos_presigned_url.get", "id"),
					resource.TestCheckResourceAttrSet("outscale_oos_presigned_url.put", "id"),
					resource.TestCheckResourceAttrSet("outscale_oos_presigned_url.head", "id"),
					resource.TestCheckResourceAttrSet("outscale_oos_presigned_url.delete", "id"),
				),
			},
		},
	})
}

func testAccPresignedURLConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  name = %q
}

resource "outscale_oos_object" "object" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "object.txt"
  content = "test object content"
}

resource "outscale_oos_presigned_url" "get" {
  bucket     = outscale_oos_bucket.bucket.id
  key        = outscale_oos_object.object.id
  method     = "GET"
  expiration = "10s"
}

resource "outscale_oos_presigned_url" "put" {
  bucket     = outscale_oos_bucket.bucket.id
  key        = outscale_oos_object.object.id
  method     = "PUT"
}

resource "outscale_oos_presigned_url" "head" {
  bucket     = outscale_oos_bucket.bucket.id
  method     = "HEAD"
}

resource "outscale_oos_presigned_url" "delete" {
  bucket     = outscale_oos_bucket.bucket.id
  method     = "DELETE"
}
`, bucketName)
}
