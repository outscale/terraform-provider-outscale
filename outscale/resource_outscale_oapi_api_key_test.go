package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/icu"
)

func TestAccOutscaleOAPIAccessKey_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var conf icu.AccessKeyMetadata

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIAccessKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIAccessKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIAccessKeyExists("outscale_oapi_api_key.a_key", &conf),
					testAccCheckOutscaleOAPIAccessKeyAttributes(&conf),
					resource.TestCheckResourceAttrSet("outscale_oapi_api_key.a_key", "secret_key"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIAccessKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).ICU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_oapi_api_key" {
			continue
		}

		// Try to get access key
		resp, err := conn.API.ListAccessKeys(nil)
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
		if err == nil {
			if len(resp.AccessKeyMetadata) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

	}

	return nil
}

func testAccCheckOutscaleOAPIAccessKeyExists(n string, res *icu.AccessKeyMetadata) resource.TestCheckFunc {
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

		if len(resp.AccessKeyMetadata) != 1 {
			return fmt.Errorf("Access Key not found not found")
		}

		*res = *resp.AccessKeyMetadata[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIAccessKeyAttributes(accessKeyMetadata *icu.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *accessKeyMetadata.Status != "Active" {
			return fmt.Errorf("Bad status: %s", *accessKeyMetadata.Status)
		}

		return nil
	}
}

const testAccOutscaleOAPIAccessKeyConfig = `
resource "outscale_oapi_api_key" "a_key" {
	#api_key_id = "AK"
	#secret_key = "SK"
}
`
