package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_policies_linked_to_user_group_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPoliciesLinkedConfig,
			},
		},
	})
}

const testAccDataPoliciesToGroupConfig = `
resource "outscale_user_group" "userGroupPolicies01" {
    user_group_name = "usergroup"
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
   policy_name = "policiesGroupLinked"
}
resource "outscale_policy" "groupLinked_02" {
   description = "Example Linked policy to group"
   document    = "{\"Statement\": [ {\"Effect\": \"Deny\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
   path        = "/Okhtgroup/"
   policy_name = "DenyPolicy"
}
data "outscale_policies_linked_to_user_group" "policiesLinkToGroup" {
	user_group_name = outscale_user_group.userGroupPolicies01.user_group_name
	depends_on = [outscale_user_group.userGroupPolicies01]
}
`
