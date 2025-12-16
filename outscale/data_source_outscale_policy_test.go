package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_data_policy_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_policy.data_test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

const testAccPolicyDataConfig = `
	resource "outscale_policy" "data_policy" {
		policy_name = "TestACC_resoucePolicy"
		document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
		path = "/"
	}
	data "outscale_policy" "data_test" {
		policy_orn = outscale_policy.data_policy.orn
	}
`
