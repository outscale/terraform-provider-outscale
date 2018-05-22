package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func TestAccOutscaleAccessKey_basic(t *testing.T) {
	var conf icu.AccessKeyMetadata

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleAccessKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists("outscale_api_key.outscale_api_key", &conf),
					testAccCheckOutscaleAccessKeyAttributes(&conf),
					resource.TestCheckResourceAttrSet("outscale_api_key.outscale_api_key", "secret_access_key"),
				),
			},
		},
	})
}

func testAccCheckOutscaleAccessKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).ICU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_api_key" {
			continue
		}

		// Try to get access key
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.ListAccessKeys(&icu.ListAccessKeysInput{})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
		if err == nil {
			return nil
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

		conn := testAccProvider.Meta().(*OutscaleClient).ICU

		resp, err := conn.API.ListAccessKeys(nil)
		if err != nil {
			return err
		}

		*res = *resp.AccessKeyMetadata[0]

		return nil
	}
}

func testAccCheckOutscaleAccessKeyAttributes(accessKeyMetadata *icu.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *accessKeyMetadata.Status != "ACTIVE" {
			return fmt.Errorf("Bad status: %s", *accessKeyMetadata.Status)
		}

		return nil
	}
}

const testAccOutscaleAccessKeyConfig = `
resource "outscale_api_key" "outscale_api_key" {
  tag = {
    Name = "test"
  }
}
`
