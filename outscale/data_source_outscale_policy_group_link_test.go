package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIPolicyGroupLinkDataSource(t *testing.T) {
	t.Skip()

	rString := acctest.RandString(8)
	groupName := fmt.Sprintf("tf-acc-group-gpa-basic-%s", rString)
	policyName := fmt.Sprintf("tf-acc-policy-gpa-basic-%s", rString)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIDSGroupPolicyAttachConfig(groupName, policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyGLDataSourceID("data.outscale_policy.policy_ds"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPolicyGLDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Can't find policy data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Policy data source ID not set")
		}

		return nil
	}
}

func testAccOutscaleOAPIDSGroupPolicyAttachConfig(groupName, policyName string) string {
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

		data "outscale_policy_group_link" "outscale_policy_group_link" {
				group_name = "${outscale_group.group.group_name}"
		}
	`, groupName, policyName)
}
