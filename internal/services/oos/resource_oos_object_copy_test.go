package oos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_ObjectCopy(t *testing.T) {
	resourceName := "outscale_oos_object_copy.copy"
	bucketName := acctest.RandomWithPrefix("testacc-oos-copy")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccObjectCopyConfig(bucketName, "private", "initial", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.copy", "initial"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
				),
			},
			{
				Config: testAccObjectCopyConfig(bucketName, "public-read", "updated", false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(resourceName, "tags.copy", "updated"),
				),
			},
			{
				Config: testAccObjectCopyConfig(bucketName, "private", "updated", true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "metadata.key", "value2"),
					resource.TestCheckResourceAttr(resourceName, "content_language", "en"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
				),
			},
		},
	})
}

func testAccObjectCopyConfig(bucketName, acl, tag string, replaceMetadata bool) string {
	metadata := "metadata_directive = \"COPY\""
	if replaceMetadata {
		metadata = `
  metadata_directive = "REPLACE"
  content_language   = "en"
  expires            = "2030-01-01T00:00:00Z"
  metadata = {
    key = "value2"
  }`
	}

	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  name = %q
}

resource "outscale_oos_object" "source" {
  bucket        = outscale_oos_bucket.bucket.id
  key           = "source"
  content       = "content"
  cache_control = "no-cache"
  content_type  = "text/plain"
  metadata = {
    key = "value1"
  }
  tags = {
    key = "value"
  }
}

resource "outscale_oos_object_copy" "copy" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "destination"
  acl    = %q
  copy_source = {
    path = "${outscale_oos_bucket.bucket.id}/${outscale_oos_object.source.key}"
  }

  %s
  tagging_directive = "REPLACE"
  tags = {
    copy = %q
  }
}
`, bucketName, acl, metadata, tag)
}
