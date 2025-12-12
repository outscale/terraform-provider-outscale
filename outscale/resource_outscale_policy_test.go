package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_policy_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_policy.basic_policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

const testAccPolicyBasicConfig = `
	resource "outscale_policy" "basic_policy" {
	  policy_name = "TestACC_resoucePolicy"
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	  path = "/"
	}`
