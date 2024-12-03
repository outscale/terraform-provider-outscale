package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_policiesLinkedToGroup_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_managed_policies_linked_to_user_group.dataPoliciesLinkedToGroup"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPoliciesLinkedGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policies.#"),
				),
			},
		},
	})
}

const testAccDataPoliciesLinkedGroupConfig = `
resource "outscale_user_group" "groupPolicies01" {
    user_group_name = "userGroupName"
    path      = "/GroupPolicies/"
    policy {
	   policy_orn = outscale_policy.GpolicyLinked_01.orn
	}
	policy {
	   policy_orn = outscale_policy.GpolicyLinked_02.orn
	}
}

resource "outscale_policy" "GpolicyLinked_01" {
   description = "Example Linked to group"
   document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
   path        = "/Allow_test/"
   policy_name = "policiesLinkedToGroup"
}
resource "outscale_policy" "GpolicyLinked_02" {
   description = "Example Linked policy to group"
   document    = "{\"Statement\": [ {\"Effect\": \"Deny\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
   path        = "/OkhtTest/"
   policy_name = "policyGroupStopAll"
}
data "outscale_managed_policies_linked_to_user_group" "dataPoliciesLinkedToGroup" {
    filter {
       name   = "path_prefix"
       values = [outscale_user_group.groupPolicies01.path]
    }
    filter {
       name   = "user_group_ids"
       values = [outscale_user_group.groupPolicies01.user_group_id]
    }
	user_group_name = outscale_user_group.groupPolicies01.user_group_name
}`
