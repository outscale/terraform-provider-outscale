package oos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_Object(t *testing.T) {
	resourceName := "outscale_oos_object.object"
	bucketResourceName := "outscale_oos_bucket.bucket"
	bucketName := acctest.RandomWithPrefix("testacc-oos-object")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccObjectConfig(bucketName, "private", "test object content", "test", "value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "key", "directory/object.txt"),
					resource.TestCheckResourceAttr(resourceName, "id", "directory/object.txt"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccObjectConfig(bucketName, "public-read", "test object content", "updated", "value"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckNoResourceAttr(resourceName, "tags.test"),
					resource.TestCheckResourceAttr(resourceName, "tags.updated", "value"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccObjectConfig(bucketName, "private", "updated object content", "updated", "value"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "content", "updated object content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "content_length", "22"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
		},
	})
}

func TestAccOOS_Object_GrantUpdate(t *testing.T) {
	resourceName := "outscale_oos_object.object"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccObjectGrantConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccObjectGrantConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "grant.read.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccObjectACLConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckNoResourceAttr(resourceName, "grant"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccObjectPrivateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "acl"),
					resource.TestCheckNoResourceAttr(resourceName, "grant"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
		},
	})
}

func TestAccOOS_Object_ContentB64(t *testing.T) {
	resourceName := "outscale_oos_object.object"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccObjectContentB64Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "content_base64"),
				),
			},
		},
	})
}

func testAccObjectConfig(bucketName, acl, content, tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  name = %q
}

resource "outscale_oos_object" "object" {
  bucket              = outscale_oos_bucket.bucket.id
  key                 = "directory/object.txt"
  content             = %q
  acl                 = %q
  cache_control       = "no-cache"
  content_disposition = "inline"
  content_encoding    = "identity"
  content_language    = "en"
  content_type        = "text/plain"
  expires             = "2030-01-01T00:00:00Z"
  metadata = {
    test = "value"
  }
  encryption_type = "AES256"
  tags = {
    %s = %q
  }
}
`, bucketName, content, acl, tagKey, tagValue)
}

func testAccObjectGrantConfig(includeRead bool) string {
	read := ""
	if includeRead {
		read = `
    read = {
      ids = [outscale_oos_bucket.bucket.permissions.owner.id]
    }`
	}

	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_object" "object" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "object.txt"
  content   = "test object content"
  grant = {
    full_control = {
      ids = [outscale_oos_bucket.bucket.permissions.owner.id]
    }%s
  }
}
`, read)
}

var testAccObjectACLConfig = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_object" "object" {
  bucket  = outscale_oos_bucket.bucket.id
  key     = "object.txt"
  content = "test object content"
  acl     = "public-read"
}
`

var testAccObjectPrivateConfig = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_object" "object" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "object.txt"
  content   = "test object content"
}
`

var testAccObjectContentB64Config = `
resource "outscale_oos_bucket" "bucket" {}

resource "outscale_oos_object" "object" {
  bucket = outscale_oos_bucket.bucket.id
  key    = "object.txt"
  content_base64   = base64encode("test object content")
}
`
