package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIPolicyDataSource_Instance(t *testing.T) {
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIPolicyDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPolicyDataSourceID("data.outscale_policy.policy_ds"),
					resource.TestCheckResourceAttr("data.outscale_policy.policy_ds", "policy_name", "test-name"),
					resource.TestCheckResourceAttr("data.outscale_policy.policy_ds", "path", "/"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPolicyDataSourceID(n string) resource.TestCheckFunc {
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

const testAccCheckOutscaleOAPIPolicyDataSourceConfig = `
	resource "outscale_policy" "policy" {
		path = "/test1"
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
		path = "${outscale_policy.policy.arn}"
	}
`
