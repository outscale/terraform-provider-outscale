package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestOutscaleOAPIPolicyVersionDataSource(t *testing.T) {
	t.Skip()

	rName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPolicyVersionDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.outscale_policy_version.policy_version_ds", "version_id"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_version.policy_version_ds", "document"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_version.policy_version_ds", "request_id"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_version.policy_version_ds", "is_default_version"),
				),
			},
		},
	})
}

func testAccOutscaleOAPIPolicyVersionDataSourceConfig(r string) string {
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

		data "outscale_policy_version" "policy_version_ds" {
			policy_arn = "${outscale_policy.outscale_policy.arn}",
			version_id = "${outscale_policy_version.policy.id}",
		}
	`, r)
}
