package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestOutscaleOAPIPolicy_namePrefix(t *testing.T) {
	var out eim.GetPolicyOutput

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPolicyPrefixNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyExists("outscale_policy.policy", &out),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPolicyDestroy(s *terraform.State) error {
	iamconn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_policy" {
			continue
		}

		_, err := iamconn.API.GetPolicy(&eim.GetPolicyInput{
			PolicyArn: aws.String(rs.Primary.ID),
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

func testAccCheckOutscaleOAPIPolicyExists(reso string, res *eim.GetPolicyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[reso]
		if !ok {
			return fmt.Errorf("Not found: %s", reso)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Policy name is set")
		}

		eimconn := testAccProvider.Meta().(*OutscaleClient).EIM

		var err error
		var resp *eim.GetPolicyOutput

		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			resp, err = eimconn.API.GetPolicy(&eim.GetPolicyInput{
				PolicyArn: aws.String(rs.Primary.Attributes["arn"]),
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

		*res = *resp

		return nil
	}
}

func testAccCheckOutscaleOAPIPolicyGeneratedNamePrefix(resource, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		name, ok := r.Primary.Attributes["name"]
		if !ok {
			return fmt.Errorf("Name attr not found: %#v", r.Primary.Attributes)
		}
		if !strings.HasPrefix(name, prefix) {
			return fmt.Errorf("Name: %q, does not have prefix: %q", name, prefix)
		}
		return nil
	}
}

const testAccOutscaleOAPIPolicyPrefixNameConfig = `
resource "outscale_policy" "policy" {
	path = "/"
  policy_name = "test-name"
  document = <<EOF
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
`
