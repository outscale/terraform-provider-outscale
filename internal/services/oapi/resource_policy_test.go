package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_Policy_Basic(t *testing.T) {
	resourceName := "outscale_policy.basic_policy"
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyBasicConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_Policy_Migration(t *testing.T) {
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.3.0", testAccPolicyBasicConfig(policyName)),
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
