package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestOutscaleOAPIPolicyDefaultVersion_namePrefix(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPolicyDefaultVersionPrefixNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyDefaultVersionExists("outscale_policy_default_version.outscale_policy_default_version"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPolicyDefaultVersionExists(reso string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[reso]
		if !ok {
			return fmt.Errorf("Not found: %s", reso)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Policy Default version is set")
		}

		return nil
	}
}

const testAccOutscaleOAPIPolicyDefaultVersionPrefixNameConfig = `
resource "outscale_policy" "outscale_policy" {
	path = "/"
  policy_name = "test-name"
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

resource "outscale_policy_default_version" "outscale_policy_default_version" {
    policy_arn = "${outscale_policy.outscale_policy.arn}"
    version_id = "v1"
}

`
