package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestOutscaleOAPIPolicyVersionsDataSource(t *testing.T) {
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
				Config: testAccOutscaleOAPIPolicyVersionsDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_policy_versions.policy_versions_ds", "versions.#", "1"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_versions.policy_versions_ds", "versions.0.version_id"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_versions.policy_versions_ds", "versions.0.document"),
					resource.TestCheckResourceAttrSet("data.outscale_policy_versions.policy_versions_ds", "versions.0.is_default_version"),
				),
			},
		},
	})
}

func testAccOutscaleOAPIPolicyVersionsDataSourceConfig(r string) string {
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

		data "outscale_policy_versions" "policy_versions_ds" {
			policy_arn = "${outscale_policy.outscale_policy.arn}"
		}
		`, r)
}
