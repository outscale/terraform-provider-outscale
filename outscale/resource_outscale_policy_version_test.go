package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestOutscalePolicyVersion(t *testing.T) {
	var out eim.GetPolicyOutput
	rName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscalePolicyVersionConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePolicyVersionExists("outscale_policy_version.policy", &out),
				),
			},
		},
	})
}

func testAccCheckOutscalePolicyVersionExists(reso string, res *eim.GetPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[reso]
		if !ok {
			return fmt.Errorf("Not found: %s", reso)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Policy name is set")
		}

		return nil
	}
}

func testAccOutscalePolicyVersionConfig(r string) string {

	return fmt.Sprintf(`
resource "outscale_policy" "outscale_policy" {
	path = "/"
  policy_name = "test-name-%s"
  policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "outscale_policy_version" "policy" {
  policy_arn = "${outscale_policy.outscale_policy.arn}"
  policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
`, r)
}
