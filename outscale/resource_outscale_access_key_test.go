package outscale

import (
	"fmt"
	"strings"
	"testing"

	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func TestAccOutscaleAccessKey_basic(t *testing.T) {
	var conf icu.AccessKeyMetadata
	rName := fmt.Sprintf("test-user-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleAccessKeyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists("outscale_api_key.a_key", &conf),
					testAccCheckOutscaleAccessKeyAttributes(&conf),
					resource.TestCheckResourceAttrSet("outscale_api_key.a_key", "secret_key_id"),
				),
			},
		},
	})
}

func testAccCheckOutscaleAccessKeyDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*OutscaleClient).ICU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_api_key" {
			continue
		}

		// Try to get access key

		var err error
		var resp *icu.ListAccessKeysOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			resp, err = iamconn.ICU.ListAccessKeys(&icu.ListAccessKeysInput{})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err == nil {
			if len(resp.AccessKeyMetadata) > 0 {
				return fmt.Errorf("still exist.")
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "NoSuchEntity" {
			return err
		}
	}

	return nil
}

func testAccCheckOutscaleAccessKeyExists(n string, res *icu.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role name is set")
		}

		iamconn := testAccProvider.Meta().(*OutscaleClient).ICU
		name := rs.Primary.Attributes["access_key_id"]

		var err error
		var resp *icu.ListAccessKeysOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			resp, err = iamconn.ICU.ListAccessKeys(&icu.ListAccessKeysInput{})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(resp.AccessKeyMetadata) != 1 ||
			*resp.AccessKeyMetadata[0].AccessKeyId != name {
			return fmt.Errorf("User not found not found")
		}

		*res = *resp.AccessKeyMetadata[0]

		return nil
	}
}

func testAccCheckOutscaleAccessKeyAttributes(accessKeyMetadata *icu.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *accessKeyMetadata.Status != "Active" {
			return fmt.Errorf("Bad status: %s", *accessKeyMetadata.Status)
		}

		return nil
	}
}

func testAccOutscaleAccessKeyConfig(rName string) string {
	return fmt.Sprint(`
resource "outscale_api_key" "a_key" {}
`)
}
