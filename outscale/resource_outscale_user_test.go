package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleUser_basic(t *testing.T) {
	var conf eim.GetUserOutput

	name1 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	name2 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	path1 := "/"
	path2 := "/path2/"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleUserConfig(name1, path1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserExists("outscale_user.user", &conf),
					testAccCheckOutscaleUserAttributes(&conf, name1, "/"),
				),
			},
			resource.TestStep{
				Config: testAccOutscaleUserConfig(name2, path2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserExists("outscale_user.user", &conf),
					testAccCheckOutscaleUserAttributes(&conf, name2, "/path2/"),
				),
			},
		},
	})
}

func testAccCheckOutscaleUserDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_user" {
			continue
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = iamconn.API.GetUser(&eim.GetUserInput{
				UserName: aws.String(rs.Primary.ID),
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
			return fmt.Errorf("still exist")
		}

		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
	}

	return nil
}

func testAccCheckOutscaleUserExists(n string, res *eim.GetUserOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User name is set")
		}

		iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

		var err error
		var resp *eim.GetUserOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = iamconn.API.GetUser(&eim.GetUserInput{
				UserName: aws.String(rs.Primary.ID),
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

		*res = *resp

		return nil
	}
}

func testAccCheckOutscaleUserAttributes(user *eim.GetUserOutput, name string, path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *user.GetUserResult.User.UserName != name {
			return fmt.Errorf("Bad name: %s", *user.GetUserResult.User.UserName)
		}

		if *user.GetUserResult.User.Path != path {
			return fmt.Errorf("Bad path: %s", *user.GetUserResult.User.Path)
		}

		return nil
	}
}

func testAccOutscaleUserConfig(r, p string) string {
	return fmt.Sprintf(`
resource "outscale_user" "user" {
	user_name = "%s"
	path = "%s"
}`, r, p)
}
