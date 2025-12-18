package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_policy_basic(t *testing.T) {
	resourceName := "outscale_policy.basic_policy"
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyBasicConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

func testAccPolicyBasicConfig(policyName string) string {
	return fmt.Sprintf(`
	resource "outscale_policy" "basic_policy" {
	  policy_name = "%s"
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	  path = "/"
	}`, policyName)
}
