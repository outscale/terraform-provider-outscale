package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleOAPIUserAPIKey_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	var conf eim.AccessKeyMetadata
	rName := fmt.Sprintf("test-user-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIUserAPIKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIUserAPIKeyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserAPIKeyExists("outscale_user_api_keys.a_key", &conf),
					testAccCheckOutscaleOAPIUserAPIKeyAttributes(&conf),
					resource.TestCheckResourceAttrSet("outscale_user_api_keys.a_key", "secret_key"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIUserAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).EIM
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_user_api_keys" {
			continue
		}

		var err error
		var resp *eim.ListAccessKeysOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.ListAccessKeys(&eim.ListAccessKeysInput{
				UserName: aws.String(rs.Primary.Attributes["user_name"]),
			})

			if err != nil {
				if strings.Contains(err.Error(), "Throttling:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err == nil {
			if len(resp.ListAccessKeysResult.AccessKeyMetadata) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		if !strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return err
		}
		return nil
	}
	return nil
}
func testAccCheckOutscaleOAPIUserAPIKeyExists(n string, res *eim.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role name is set")
		}
		conn := testAccProvider.Meta().(*OutscaleClient).EIM
		name := rs.Primary.Attributes["user_name"]
		var err error
		var resp *eim.ListAccessKeysOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.ListAccessKeys(&eim.ListAccessKeysInput{
				UserName: aws.String(name),
			})

			if err != nil {
				if strings.Contains(err.Error(), "Throttling:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		if len(resp.ListAccessKeysResult.AccessKeyMetadata) != 1 ||
			*resp.ListAccessKeysResult.AccessKeyMetadata[0].UserName != name {
			return fmt.Errorf("User not found not found")
		}
		*res = *resp.ListAccessKeysResult.AccessKeyMetadata[0]
		return nil
	}
}
func testAccCheckOutscaleOAPIUserAPIKeyAttributes(accessKeyMetadata *eim.AccessKeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(*accessKeyMetadata.UserName, "test-user") {
			return fmt.Errorf("Bad username: %s", *accessKeyMetadata.UserName)
		}
		if *accessKeyMetadata.Status != "Active" {
			return fmt.Errorf("Bad status: %s", *accessKeyMetadata.Status)
		}
		return nil
	}
}

func testAccOutscaleOAPIUserAPIKeyConfig(rName string) string {
	return fmt.Sprintf(`
resource "outscale_user" "a_user" {
        user_name = "%s"
}
resource "outscale_user_api_keys" "a_key" {
        user_name = "${outscale_user.a_user.user_name}"
}
`, rName)
}
