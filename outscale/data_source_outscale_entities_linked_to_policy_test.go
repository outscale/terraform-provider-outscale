package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_entities_linked_to_policy_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_entities_linked_to_policy.entitiesLinked"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEntitiesLinkedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_entities.#"),
				),
			},
		},
	})
}

const testAccDataEntitiesLinkedConfig = `
resource "outscale_user" "user_01" {
    user_name = "userLedGroup"
    path      = "/linkedUser/"
    policy {
	   policy_orn = outscale_policy.policyEntities_01.orn
	}
}
resource "outscale_user_group" "uGroupLinked" {
	user_group_name = "GLinkedTestACC"
	path = "/"
	user {
        user_name = outscale_user.user_01.user_name
        path      = "/linkedUser/"
    }
	policy {
	   policy_orn = outscale_policy.policyEntities_01.orn
	}
	depends_on = [outscale_user.user_01]
}
resource "outscale_user_group" "GroupLinkedPolicy" {
	user_group_name = "GroupPolicyTestAcc"
	path = "/TestPath/"
	policy {
	   policy_orn = outscale_policy.policyEntities_01.orn
	}
}
resource "outscale_policy" "policyEntities_01" {
   description = "Example Entities Linked to policy"
   document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
   path        = "/Okht_test/"
   policy_name = "policyEntitiesLinked"
}

data "outscale_entities_linked_to_policy" "entitiesLinked" {
	policy_orn = outscale_policy.policyEntities_01.orn
	depends_on = [outscale_user_group.uGroupLinked, outscale_user_group.GroupLinkedPolicy, outscale_user.user_01]
}`
