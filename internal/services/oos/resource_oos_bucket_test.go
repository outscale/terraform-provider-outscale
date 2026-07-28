package oos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOOS_Bucket_Basic(t *testing.T) {
	resourceName := "outscale_oos_bucket.bucket"
	name := acctest.RandomWithPrefix("testacc-oos-bucket")

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "force_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "object_lock.enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "permissions.owner.id"),
				),
			},
		},
	})
}

func TestAccOOS_Bucket_ACLUpdate(t *testing.T) {
	resourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketACLConfig("private"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccOOSBucketACLConfig("public-read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccOOSBucketACLConfig("private"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
		},
	})
}

func TestAccOOS_Bucket_GrantUpdate(t *testing.T) {
	resourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketGrantConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccOOSBucketGrantConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "grant.read.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccOOSBucketGrantConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "grant.read"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccOOSBucketGrantACLConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckNoResourceAttr(resourceName, "grant"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "2"),
				),
			},
			{
				Config: testAccOOSBucketGrantConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "acl"),
					resource.TestCheckResourceAttr(resourceName, "grant.full_control.ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
			{
				Config: testAccOOSBucketGrantPrivateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "acl"),
					resource.TestCheckNoResourceAttr(resourceName, "grant"),
					resource.TestCheckResourceAttr(resourceName, "permissions.grants.#", "1"),
				),
			},
		},
	})
}

func TestAccOOS_Bucket_ObjectLockUpdate(t *testing.T) {
	resourceName := "outscale_oos_bucket.bucket"

	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccOOSBucketObjectLockConfig(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_lock.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock.default_retention.mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock.default_retention.days", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "object_lock.default_retention.years"),
				),
			},
			{
				Config: testAccOOSBucketObjectLockConfig(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_lock.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock.default_retention.mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock.default_retention.days", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "object_lock.default_retention.years"),
				),
			},
			{
				Config: testAccOOSBucketObjectLockWithoutRetentionConfig(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_lock.enabled", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "object_lock.default_retention"),
				),
			},
			{
				Config: testAccOOSBucketWithoutObjectLockConfig(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_lock.enabled", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "object_lock.default_retention"),
				),
			},
		},
	})
}

func testAccOOSBucketConfig(name string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  name = %q
}
`, name)
}

func testAccOOSBucketACLConfig(acl string) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  acl = %q
}
`, acl)
}

func testAccOOSBucketGrantConfig(includeRead bool) string {
	read := ""
	if includeRead {
		read = `
    read = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }`
	}

	return fmt.Sprintf(`
resource "outscale_oos_bucket" "owner" {}

resource "outscale_oos_bucket" "bucket" {
  grant = {
    full_control = {
      ids = [outscale_oos_bucket.owner.permissions.owner.id]
    }%s
  }
}
`, read)
}

func testAccOOSBucketGrantACLConfig() string {
	return `
resource "outscale_oos_bucket" "owner" {}

resource "outscale_oos_bucket" "bucket" {
  acl = "public-read"
}
`
}

func testAccOOSBucketGrantPrivateConfig() string {
	return `
resource "outscale_oos_bucket" "owner" {}

resource "outscale_oos_bucket" "bucket" {}
`
}

func testAccOOSBucketObjectLockConfig(days int) string {
	return fmt.Sprintf(`
resource "outscale_oos_bucket" "bucket" {
  object_lock = {
    enabled = true
    default_retention = {
      days = %d
    }
  }
}
`, days)
}

func testAccOOSBucketObjectLockWithoutRetentionConfig() string {
	return `
resource "outscale_oos_bucket" "bucket" {
  object_lock = {
    enabled = true
  }
}
`
}

func testAccOOSBucketWithoutObjectLockConfig() string {
	return `
resource "outscale_oos_bucket" "bucket" {}
`
}
