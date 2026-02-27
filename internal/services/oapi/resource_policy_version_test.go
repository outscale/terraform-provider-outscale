package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_PolicyVersion_Basic(t *testing.T) {
	resourceName := "outscale_policy_version.policy_version"
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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

func TestAccOthers_PolicyVersion_Migration(t *testing.T) {
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.3.1", testAccPolicyVersionConfig(policyName)),
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
