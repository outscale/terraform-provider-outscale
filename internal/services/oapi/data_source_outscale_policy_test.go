package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_data_policy_basic(t *testing.T) {
	resourceName := "data.outscale_policy.data_test"
	policyName := acctest.RandomWithPrefix("test-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyDataConfig(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
		},
	})
}

func testAccPolicyDataConfig(policyName string) string {
	return fmt.Sprintf(`
	resource "outscale_policy" "data_policy" {
		policy_name = "%s"
		document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
		path = "/"
	}
	data "outscale_policy" "data_test" {
		policy_orn = outscale_policy.data_policy.orn
	}
	`, policyName)
}
