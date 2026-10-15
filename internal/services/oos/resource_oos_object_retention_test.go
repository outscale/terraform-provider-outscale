package oos_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_ObjectRetention(t *testing.T) {
	resourceName := "outscale_oos_object_retention.retention"
	remainUntil := time.Now().UTC().Add(5 * time.Second)

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccObjectRetentionConfig(remainUntil.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(resourceName, "retain_until_date", remainUntil.Format(time.RFC3339)),
				),
			},
		},
	})
}

func testAccObjectRetentionConfig(retainUntilDate string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  force_delete = true
  object_lock = {
    enabled = true
  }
}

resource "outscale_oos_object" "object" {
  bucket       = outscale_oos_bucket.bucket.id
  key          = "object"
  content      = "content"
  force_delete = true
}

resource "outscale_oos_object_retention" "retention" {
  bucket            = outscale_oos_object.object.bucket
  key               = outscale_oos_object.object.key
  retain_until_date = "%s"
}

 resource "terraform_data" "wait_for_retention" {
   depends_on = [outscale_oos_object_retention.retention]

   provisioner "local-exec" {
     command = "sleep 6"
   }
}
`, retainUntilDate)
}
