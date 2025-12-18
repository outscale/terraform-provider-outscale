package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_policy_Version_basic(t *testing.T) {
	resourceName := "outscale_policy_version.policy_version"
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyVersionConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "body"),
					resource.TestCheckResourceAttr(resourceName, "version_id", "v2"),
				),
			},
		},
	})
}

func testAccPolicyVersionConfig(name string) string {
	return fmt.Sprintf(`
	resource "outscale_policy" "vers_policy" {
	  policy_name = "%s"
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	  path = "/"
	}
	resource "outscale_policy_version" "policy_version" {
	  policy_orn = outscale_policy.vers_policy.orn
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	}
`, name)
}
