package outscale

import (
	"errors"
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

func TestAccOutscaleOAPIPolicyGroup_basic(t *testing.T) {
	var groupPolicy1, groupPolicy2 eim.GetGroupPolicyOutput
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEIMGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIMGroupPolicyConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIMGroupPolicyExists(
						"outscale_group.group",
						"outscale_policy_group.foo",
						&groupPolicy1,
					),
				),
			},
			{
				Config: testAccEIMGroupPolicyConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIMGroupPolicyExists(
						"outscale_group.group",
						"outscale_policy_group.bar",
						&groupPolicy2,
					),
					testAccCheckOutscaleOAPIPolicyGroupNameChanged(&groupPolicy1, &groupPolicy2),
				),
			},
		},
	})
}

func testAccCheckEIMGroupPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_policy_group" {
			continue
		}

		group, name := resourceOutscaleOAPIPolicyGroupParseID(rs.Primary.ID)

		request := &eim.GetGroupPolicyInput{
			PolicyName: aws.String(name),
			GroupName:  aws.String(group),
		}

		var err error
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			_, err = conn.API.GetGroupPolicy(request)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Throttling") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
				continue
			}
			return err
		}

		return fmt.Errorf("still exists")
	}

	return nil
}

func testAccCheckEIMGroupPolicyExists(
	iamGroupResource string,
	iamGroupPolicyResource string,
	groupPolicy *eim.GetGroupPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[iamGroupResource]
		if !ok {
			return fmt.Errorf("Not Found: %s", iamGroupResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		policy, ok := s.RootModule().Resources[iamGroupPolicyResource]
		if !ok {
			return fmt.Errorf("Not Found: %s", iamGroupPolicyResource)
		}

		iamconn := testAccProvider.Meta().(*OutscaleClient).EIM
		group, name := resourceOutscaleOAPIPolicyGroupParseID(policy.Primary.ID)

		var err error
		var output *eim.GetGroupPolicyOutput
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			output, err = iamconn.API.GetGroupPolicy(&eim.GetGroupPolicyInput{
				GroupName:  aws.String(group),
				PolicyName: aws.String(name),
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "Throttling") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		*groupPolicy = *output

		return nil
	}
}

func testAccCheckOutscaleOAPIPolicyGroupNameChanged(i, j *eim.GetGroupPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.GetGroupPolicyResult.PolicyName) == aws.StringValue(j.GetGroupPolicyResult.PolicyName) {
			return errors.New("EIM Group Policy name did not change")
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIPolicyGroupNameMatches(i, j *eim.GetGroupPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.GetGroupPolicyResult.PolicyName) != aws.StringValue(j.GetGroupPolicyResult.PolicyName) {
			return errors.New("EIM Group Policy name did not match")
		}

		return nil
	}
}

func testAccEIMGroupPolicyConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_group" "group" {
		group_name = "test_group_%d"
		path = "/"
	}

	resource "outscale_policy_group" "foo" {
		policy_id = "foo_policy_%d"
		group_name = "${outscale_group.group.group_name}"
		policy_document = <<EOF
{
	"Version": "2012-10-17",
	"Statement": {
		"Effect": "Allow",
		"Action": "*",
		"Resource": "*"
	}
}
EOF
	}`, rInt, rInt)
}

func testAccEIMGroupPolicyConfigUpdate(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_group" "group" {
		group_name = "test_group_%d"
		path = "/"
	}

	resource "outscale_policy_group" "foo" {
		policy_id = "foo_policy_%d"
		group_name = "${outscale_group.group.group_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}

	resource "outscale_policy_group" "bar" {
		policy_id = "bar_policy_%d"
		group_name = "${outscale_group.group.group_name}"
		policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"
	}`, rInt, rInt, rInt)
}
