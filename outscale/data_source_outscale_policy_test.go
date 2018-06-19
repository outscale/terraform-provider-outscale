package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscalePolicyDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscalePolicyDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePolicyDataSourceID("data.outscale_policy.policy_ds"),
					resource.TestCheckResourceAttr("data.outscale_policy.policy_ds", "policy_name", "test-name"),
					resource.TestCheckResourceAttr("data.outscale_policy.policy_ds", "path", "/"),
				),
			},
		},
	})
}

func testAccCheckOutscalePolicyDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
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

const testAccCheckOutscalePolicyDataSourceConfig = `
resource "outscale_policy" "policy" {
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

data "outscale_policy" "policy_ds" {
	policy_arn = "${outscale_policy.policy.arn}"
}
`
