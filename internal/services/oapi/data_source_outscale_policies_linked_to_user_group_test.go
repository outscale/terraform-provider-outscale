package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_policies_linked_to_user_group_basic(t *testing.T) {
	name1 := acctest.RandomWithPrefix("testacc-policy")
	name2 := acctest.RandomWithPrefix("testacc-policy")
	groupName := acctest.RandomWithPrefix("testacc-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataPoliciesToGroupConfig(name1, name2, groupName),
			},
		},
	})
}

func testAccDataPoliciesToGroupConfig(name1, name2, groupName string) string {
	return fmt.Sprintf(`
	resource "outscale_user_group" "userGroupPolicies01" {
	    user_group_name = "%[3]s"
	    path      = "/policiesGroup/"
	    policy {
		   policy_orn = outscale_policy.groupLinked_01.orn
		}
		policy {
		   policy_orn = outscale_policy.groupLinked_02.orn
		}
	}

	resource "outscale_policy" "groupLinked_01" {
	   description = "Example Linked to user_group"
	   document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	   path        = "/TestAllOK/"
	   policy_name = "%[1]s"
	}
	resource "outscale_policy" "groupLinked_02" {
	   description = "Example Linked policy to group"
	   document    = "{\"Statement\": [ {\"Effect\": \"Deny\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	   path        = "/Okhtgroup/"
	   policy_name = "%[2]s"
	}
	data "outscale_policies_linked_to_user_group" "policiesLinkToGroup" {
		user_group_name = outscale_user_group.userGroupPolicies01.user_group_name
		depends_on = [outscale_user_group.userGroupPolicies01]
	}
	`, name1, name2, groupName)
}
