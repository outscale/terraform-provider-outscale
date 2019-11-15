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

func TestAccOutscaleOAPIUser_basic(t *testing.T) {
	t.Skip()

	var conf eim.GetUserOutput

	name1 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	name2 := fmt.Sprintf("test-user-%d", acctest.RandInt())
	path1 := "/"
	path2 := "/path2/"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIUserConfig(name1, path1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserExists("outscale_user.user", &conf),
					testAccCheckOutscaleOAPIUserAttributes(&conf, name1, "/"),
				),
			},
			resource.TestStep{
				Config: testAccOutscaleOAPIUserConfig(name2, path2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserExists("outscale_user.user", &conf),
					testAccCheckOutscaleOAPIUserAttributes(&conf, name2, "/path2/"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIUserDestroy(s *terraform.State) error {
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
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
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

func testAccCheckOutscaleOAPIUserExists(n string, res *eim.GetUserOutput) resource.TestCheckFunc {
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

		*res = *resp

		return nil
	}
}

func testAccCheckOutscaleOAPIUserAttributes(user *eim.GetUserOutput, name string, path string) resource.TestCheckFunc {
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

func testAccOutscaleOAPIUserConfig(r, p string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "user" {
			user_name = "%s"
			path      = "%s"
		}
	`, r, p)
}
