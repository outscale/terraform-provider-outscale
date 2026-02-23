package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_policies_linked_to_user_basic(t *testing.T) {
	resourceName := "data.outscale_policies_linked_to_user.policiesLinkedToUser"
	name1 := acctest.RandomWithPrefix("test-policy")
	name2 := acctest.RandomWithPrefix("test-policy")
	userName := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataPoliciesLinkedConfig(name1, name2, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policies.#"),
				),
			},
		},
	})
}

func testAccDataPoliciesLinkedConfig(name1, name2, userName string) string {
	return fmt.Sprintf(`
	resource "outscale_user" "userPolicies01" {
	    user_name = "%[3]s"
	    path      = "/policiesUser/"
	    policy {
		   policy_orn = outscale_policy.policyLinked_01.orn
		}
		policy {
		   policy_orn = outscale_policy.policyLinked_02.orn
		}
	}

	resource "outscale_policy" "policyLinked_01" {
	   description = "Example Linked to user"
	   document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	   path        = "/Allow_est/"
	   policy_name = "%[1]s"
	}
	resource "outscale_policy" "policyLinked_02" {
	   description = "Example Linked policy to user"
	   document    = "{\"Statement\": [ {\"Effect\": \"Deny\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	   path        = "/OkhtTest/"
	   policy_name = "%[2]s"
	}
	data "outscale_policies_linked_to_user" "policiesLinkedToUser" {
		user_name = outscale_user.userPolicies01.user_name
		depends_on = [outscale_user.userPolicies01]
	}
	`, name1, name2, userName)
}
