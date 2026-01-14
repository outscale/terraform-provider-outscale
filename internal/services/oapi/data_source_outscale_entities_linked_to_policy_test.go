package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_entities_linked_to_policy_basic(t *testing.T) {
	resourceName := "data.outscale_entities_linked_to_policy.entitiesLinked"

	userName := acctest.RandomWithPrefix("testacc-user")
	policyName := acctest.RandomWithPrefix("test-policy")
	groupName1 := acctest.RandomWithPrefix("testacc-usergroupname")
	groupName2 := acctest.RandomWithPrefix("testacc-usergroupname")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataEntitiesLinkedConfig(policyName, groupName1, groupName2, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy_entities.#"),
				),
			},
		},
	})
}

func testAccDataEntitiesLinkedConfig(policyName, groupName1, groupName2, userName string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "user_01" {
		    user_name = "%[4]s"
		    path      = "/linkedUser/"
		    policy {
			   policy_orn = outscale_policy.policyEntities_01.orn
			}
		}
		resource "outscale_user_group" "uGroupLinked" {
			user_group_name = "%[2]s"
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
			user_group_name = "%[3]s"
			path = "/TestPath/"
			policy {
			   policy_orn = outscale_policy.policyEntities_01.orn
			}
		}
		resource "outscale_policy" "policyEntities_01" {
		   description = "Example Entities Linked to policy"
		   document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
		   path        = "/Okht_test/"
		   policy_name = "%[1]s"
		}

		data "outscale_entities_linked_to_policy" "entitiesLinked" {
			policy_orn = outscale_policy.policyEntities_01.orn
			depends_on = [outscale_user_group.uGroupLinked, outscale_user_group.GroupLinkedPolicy, outscale_user.user_01]
		}
	`, policyName, groupName1, groupName2, userName)
}
