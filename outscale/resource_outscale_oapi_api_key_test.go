package outscale

import (
	"fmt"
	"strings"
	"testing"

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
					testAccCheckOutscaleAccessKeyExists("outscale_oapi_api_key.a_key", &conf),
					testAccCheckOutscaleAccessKeyAttributes(&conf),
					resource.TestCheckResourceAttrSet("outscale_oapi_api_key.a_key", "secret_key"),
				),
			},
		},
	})
}

func testAccCheckOutscaleAccessKeyDestroy(s *terraform.State) error {
	client_icu := testAccProvider.Meta().(*OutscaleClient).ICU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_oapi_api_key" {
			continue
		}

		// Try to get access key
		resp, err := client_icu.API.ListAccessKeys(nil)
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
		if err == nil {
			if len(resp.AccessKeyMetadata) > 0 {
				return fmt.Errorf("still exist.")
			}
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

		client_icu := testAccProvider.Meta().(*OutscaleClient).ICU

		resp, err := client_icu.API.ListAccessKeys(nil)
		if err != nil {
			return err
		}

		if len(resp.AccessKeyMetadata) != 1 {
			return fmt.Errorf("Access Key not found not found")
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

const testAccOutscaleAccessKeyConfig = `
resource "outscale_oapi_api_key" "a_key" {
	#api_key_id = "7E4U4AQ0CGLTWB78Q38V"
	#secret_key = "TDKLDVCNFDWFT6CVYBM9OPQ5YO9ZAJBN0JBJS99K"
}
`
