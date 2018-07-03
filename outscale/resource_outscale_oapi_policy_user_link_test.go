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

func TestAccOutscaleOAPIPolicyUserLink_basic(t *testing.T) {
	var out eim.ListAttachedUserPoliciesOutput
	rName := acctest.RandString(10)
	policyName1 := fmt.Sprintf("test-policy-%s", acctest.RandString(10))
	policyName2 := fmt.Sprintf("test-policy-%s", acctest.RandString(10))
	policyName3 := fmt.Sprintf("test-policy-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPolicyUserLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIUserPolicyAttachConfig(rName, policyName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyUserLinkExists("outscale_policy_user_link.test-attach", 1, &out),
					testAccCheckOutscaleOAPIPolicyUserLinkAttributes([]string{policyName1}, &out),
				),
			},
			{
				Config: testAccOutscaleOAPIUserPolicyAttachConfigUpdate(rName, policyName1, policyName2, policyName3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyUserLinkExists("outscale_policy_user_link.test-attach", 2, &out),
					testAccCheckOutscaleOAPIPolicyUserLinkAttributes([]string{policyName2, policyName3}, &out),
				),
			},
		},
	})
}
func testAccCheckOutscaleOAPIPolicyUserLinkDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckOutscaleOAPIPolicyUserLinkExists(n string, c int, out *eim.ListAttachedUserPoliciesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No policy name is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).EIM
		user := rs.Primary.Attributes["user_name"]

		var err error
		var attachedPolicies *eim.ListAttachedUserPoliciesOutput
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			attachedPolicies, err = conn.API.ListAttachedUserPolicies(&eim.ListAttachedUserPoliciesInput{
				UserName: aws.String(user),
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
			return fmt.Errorf("Error: Failed to get attached policies for user %s (%s)", user, n)
		}
		if c != len(attachedPolicies.AttachedPolicies) {
			return fmt.Errorf("Error: User (%s) has wrong number of policies attached on initial creation", n)
		}

		*out = *attachedPolicies
		return nil
	}
}
func testAccCheckOutscaleOAPIPolicyUserLinkAttributes(policies []string, out *eim.ListAttachedUserPoliciesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		matched := 0

		for _, p := range policies {
			for _, ap := range out.AttachedPolicies {
				// *ap.PolicyArn like arn:aws:eim::111111111111:policy/test-policy
				parts := strings.Split(*ap.PolicyArn, "/")
				if len(parts) == 2 && p == parts[1] {
					matched++
				}
			}
		}
		if matched != len(policies) || matched != len(out.AttachedPolicies) {
			return fmt.Errorf("Error: Number of attached policies was incorrect: expected %d matched policies, matched %d of %d", len(policies), matched, len(out.AttachedPolicies))
		}
		return nil
	}
}

func testAccOutscaleOAPIUserPolicyAttachConfig(rName, policyName string) string {
	return fmt.Sprintf(`
resource "outscale_user" "user" {
    user_name = "test-user-%s"
}

resource "outscale_policy" "policy" {
    policy_name = "%s"
    #description = "A test policy"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "eim:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_user_link" "test-attach" {
    user_name = "${outscale_user.user.user_name}"
    policy_arn = "${outscale_policy.policy.arn}"
}`, rName, policyName)
}

func testAccOutscaleOAPIUserPolicyAttachConfigUpdate(rName, policyName1, policyName2, policyName3 string) string {
	return fmt.Sprintf(`
resource "outscale_user" "user" {
    user_name = "test-user-%s"
}

resource "outscale_policy" "policy" {
    policy_name = "%s"
    #description = "A test policy"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "eim:ChangePassword"
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
    #description = "A test policy"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "eim:ChangePassword"
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
    #description = "A test policy"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "eim:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_user_link" "test-attach" {
    user_name = "${outscale_user.user.user_name}"
    policy_arn = "${outscale_policy.policy2.arn}"
}

resource "outscale_policy_user_link" "test-attach2" {
    user_name = "${outscale_user.user.user_name}"
    policy_arn = "${outscale_policy.policy3.arn}"
}`, rName, policyName1, policyName2, policyName3)
}
