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

func TestAccOutscaleEIMGroupPolicyAttachment_basic(t *testing.T) {
	var out eim.ListAttachedGroupPoliciesOutput

	rString := acctest.RandString(8)
	groupName := fmt.Sprintf("tf-acc-group-gpa-basic-%s", rString)
	policyName := fmt.Sprintf("tf-acc-policy-gpa-basic-%s", rString)
	policyName2 := fmt.Sprintf("tf-acc-policy-gpa-basic-2-%s", rString)
	policyName3 := fmt.Sprintf("tf-acc-policy-gpa-basic-3-%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleGroupPolicyAttachConfig(groupName, policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleGroupPolicyAttachmentExists("outscale_policy_group_link.test-attach", 1, &out),
					testAccCheckOutscaleGroupPolicyAttachmentAttributes([]string{policyName}, &out),
				),
			},
			resource.TestStep{
				Config: testAccOutscaleGroupPolicyAttachConfigUpdate(groupName, policyName, policyName2, policyName3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleGroupPolicyAttachmentExists("outscale_policy_group_link.test-attach", 2, &out),
					testAccCheckOutscaleGroupPolicyAttachmentAttributes([]string{policyName2, policyName3}, &out),
				),
			},
		},
	})
}
func testAccCheckOutscaleGroupPolicyAttachmentDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckOutscaleGroupPolicyAttachmentExists(n string, c int, out *eim.ListAttachedGroupPoliciesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No policy name is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).EIM
		group := rs.Primary.Attributes["group_name"]

		var err error
		var attachedPolicies *eim.ListAttachedGroupPoliciesOutput
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			attachedPolicies, err = conn.API.ListAttachedGroupPolicies(&eim.ListAttachedGroupPoliciesInput{
				GroupName: aws.String(group),
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
			return fmt.Errorf("Error: Failed to get attached policies for group %s (%s)", group, n)
		}
		if c != len(attachedPolicies.ListAttachedGroupPoliciesResult.AttachedPolicies) {
			return fmt.Errorf("Error: Group (%s) has wrong number of policies attached on initial creation", n)
		}

		*out = *attachedPolicies
		return nil
	}
}
func testAccCheckOutscaleGroupPolicyAttachmentAttributes(policies []string, out *eim.ListAttachedGroupPoliciesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		matched := 0

		for _, p := range policies {
			for _, ap := range out.ListAttachedGroupPoliciesResult.AttachedPolicies {
				// *ap.PolicyArn like arn:aws:iam::111111111111:policy/test-policy
				parts := strings.Split(*ap.PolicyArn, "/")
				if len(parts) == 2 && p == parts[1] {
					matched++
				}
			}
		}
		if matched != len(policies) || matched != len(out.ListAttachedGroupPoliciesResult.AttachedPolicies) {
			return fmt.Errorf("Error: Number of attached policies was incorrect: expected %d matched policies, matched %d of %d", len(policies), matched, len(out.ListAttachedGroupPoliciesResult.AttachedPolicies))
		}
		return nil
	}
}

func testAccOutscaleGroupPolicyAttachConfig(groupName, policyName string) string {
	return fmt.Sprintf(`
resource "outscale_group" "group" {
    group_name = "%s"
}

resource "outscale_policy" "policy" {
    policy_name = "%s"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_group_link" "test-attach" {
    group_name = "${outscale_group.group.group_name}"
    policy_arn = "${outscale_policy.policy.arn}"
}
`, groupName, policyName)
}

func testAccOutscaleGroupPolicyAttachConfigUpdate(groupName, policyName, policyName2, policyName3 string) string {
	return fmt.Sprintf(`
resource "outscale_group" "group" {
    group_name = "%s"
}

resource "outscale_policy" "policy" {
    policy_name = "%s"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy" "policy2" {
    policy_name = "%s"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy" "policy3" {
    policy_name = "%s"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_group_link" "test-attach" {
    group_name = "${outscale_group.group.group_name}"
    policy_arn = "${outscale_policy.policy2.arn}"
}

resource "outscale_policy_group_link" "test-attach2" {
    group_name = "${outscale_group.group.group_name}"
    policy_arn = "${outscale_policy.policy3.arn}"
}
`, groupName, policyName, policyName2, policyName3)
}
