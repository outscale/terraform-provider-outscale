package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_policy_Version_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_policy_version.policy_version"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyVersionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "body"),
					resource.TestCheckResourceAttr(resourceName, "version_id", "v2"),
				),
			},
		},
	})
}

const testAccPolicyVersionConfig = `
	resource "outscale_policy" "vers_policy" {
	  policy_name = "TestACC_VersionPolicy"
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	  path = "/"
	}
	resource "outscale_policy_version" "policy_version" {
	  policy_orn = outscale_policy.vers_policy.orn
	  document = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
	}
`
