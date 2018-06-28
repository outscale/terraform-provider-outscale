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

func TestAccOutscaleUserPolicy_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserPolicyConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserPolicy(
						"outscale_user.user",
						"outscale_policy_user.foo",
					),
				),
			},
			{
				Config: testAccOutscaleUserPolicyConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserPolicy(
						"outscale_user.user",
						"outscale_policy_user.bar",
					),
				),
			},
		},
	})
}

func TestAccOutscaleUserPolicy_namePrefix(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_policy_user.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserPolicyConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserPolicy(
						"outscale_user.test",
						"outscale_policy_user.test",
					),
				),
			},
		},
	})
}

func TestAccOutscaleUserPolicy_generatedName(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_policy_user.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleUserPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserPolicyConfig3(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleUserPolicy(
						"outscale_user.test",
						"outscale_policy_user.test",
					),
				),
			},
		},
	})
}

func testAccCheckOutscaleUserPolicyDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_policy_user" {
			continue
		}

		role, name := resourceOutscaleUserPolicyParseID(rs.Primary.ID)

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

func testAccCheckOutscaleUserPolicy(
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
		username, name := resourceOutscaleUserPolicyParseID(policy.Primary.ID)
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

func testAccOutscaleUserPolicyConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "user" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "foo" {
		policy_name = "foo_policy_%d"
		user_name = "${outscale_user.user.user_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"

		depends_on = ["outscale_user.user"]
	}`, rInt, rInt)
}

func testAccOutscaleUserPolicyConfig2(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "test" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "test" {
		user_name = "${outscale_user.test.user_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}`, rInt)
}

func testAccOutscaleUserPolicyConfig3(rInt int) string {
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

func testAccOutscaleUserPolicyConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_user" "user" {
		user_name = "test_user_%d"
		path = "/"
	}

	resource "outscale_policy_user" "foo" {
		policy_name = "foo_policy_%d"
		user_name = "${outscale_user.user.user_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}

	resource "outscale_policy_user" "bar" {
		policy_name = "bar_policy_%d"
		user_name = "${outscale_user.user.user_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}`, rInt, rInt, rInt)
}
