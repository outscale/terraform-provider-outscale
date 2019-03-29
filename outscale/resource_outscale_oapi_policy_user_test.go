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

func TestAccOutscaleOAPIUserPolicy_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIUserPolicyConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserPolicy(
						"outscale_user.user",
						"outscale_policy_user.foo",
					),
				),
			},
			{
				Config: testAccOutscaleOAPIUserPolicyConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserPolicy(
						"outscale_user.user",
						"outscale_policy_user.bar",
					),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIUserPolicy_namePrefix(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_policy_user.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIUserPolicyConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserPolicy(
						"outscale_user.test",
						"outscale_policy_user.test",
					),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIUserPolicy_generatedName(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_policy_user.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIUserPolicyConfig3(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIUserPolicy(
						"outscale_user.test",
						"outscale_policy_user.test",
					),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIUserPolicyDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_policy_user" {
			continue
		}

		role, name := resourceOutscaleOAPIUserPolicyParseID(rs.Primary.ID)

		request := &eim.GetRolePolicyInput{
			PolicyName: aws.String(name),
			RoleName:   aws.String(role),
		}

		var err error
		var getResp *eim.GetRolePolicyOutput
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			getResp, err = iamconn.API.GetRolePolicy(request)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
				return nil
			}
			return fmt.Errorf("Error reading Outscale policy %s from role %s: %s", name, role, err)
		}

		if getResp != nil {
			return fmt.Errorf("Found Outscale Role, expected none: %v", getResp)
		}
	}

	return nil
}

func testAccCheckOutscaleOAPIUserPolicy(
	iamUserResource string,
	iamUserPolicyResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[iamUserResource]
		if !ok {
			return fmt.Errorf("Not Found: %s", iamUserResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		policy, ok := s.RootModule().Resources[iamUserPolicyResource]
		if !ok {
			return fmt.Errorf("Not Found: %s", iamUserPolicyResource)
		}

		iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

		var err error
		username, name := resourceOutscaleOAPIUserPolicyParseID(policy.Primary.ID)
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			_, err = iamconn.API.GetUserPolicy(&eim.GetUserPolicyInput{
				UserName:   aws.String(username),
				PolicyName: aws.String(name),
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccOutscaleOAPIUserPolicyConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "user" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "foo" {
		policy_id = "foo_policy_%d"
		user_name = "${outscale_user.user.user_name}"

		depends_on = ["outscale_user.user"]
	}`, rInt, rInt)
}

func testAccOutscaleOAPIUserPolicyConfig2(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "test" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "test" {
		user_name = "${outscale_user.test.user_name}"
	}`, rInt)
}

func testAccOutscaleOAPIUserPolicyConfig3(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "test" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "test" {
		user_name = "${outscale_user.test.user_name}"
		policy_docuemnt = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}`, rInt)
}

func testAccOutscaleOAPIUserPolicyConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "user" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "foo" {
		policy_id = "foo_policy_%d"
		user_name = "${outscale_user.user.user_name}"
	}

	resource "outscale_policy_user" "bar" {
		policy_id = "bar_policy_%d"
		user_name = "${outscale_user.user.user_name}"
	}`, rInt, rInt, rInt)
}
