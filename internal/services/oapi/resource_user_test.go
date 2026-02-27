package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_User_Basic(t *testing.T) {
	resourceName := "outscale_user.basic_user"
	userName := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserBasicConfig(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_User_Policy(t *testing.T) {
	resourceName := "outscale_user.user_policy"
	userName := acctest.RandomWithPrefix("testacc-user")
	policyName := acctest.RandomWithPrefix("testacc-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserWithPolicy(policyName, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.#"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_User_Update(t *testing.T) {
	resourceName := "outscale_user.update_user"
	userName := acctest.RandomWithPrefix("testacc-user")
	userName2 := acctest.RandomWithPrefix("testacc-user")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleUserUpdatedConfig(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
				),
			},
			{
				Config: testAccOutscaleUserUpdatedConfig(userName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName2),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_User_UppercaseVersionID(t *testing.T) {
	resourceName := "outscale_user.user_policy_version"
	userName := acctest.RandomWithPrefix("testacc-user")
	policyName := acctest.RandomWithPrefix("testacc-policy")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserWithPolicyVersionUpper(policyName, userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttrSet(resourceName, "policy.#"),
					resource.TestCheckResourceAttr(resourceName, "policy.0.default_version_id", "V2"),
				),
			},
		},
	})
}

func TestAccOthers_User_Migration(t *testing.T) {
	userName := acctest.RandomWithPrefix("testacc-user")
	policyName := acctest.RandomWithPrefix("testacc-policy")

	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.3.1", testAccOutscaleUserBasicConfig(userName), testAccUserWithPolicy(policyName, userName)),
	})
}

func testAccOutscaleUserBasicConfig(userName string) string {
	return fmt.Sprintf(`
    resource "outscale_access_key" "access_key01" {
		state       = "ACTIVE"
		user_name   = outscale_user.basic_user.user_name
		depends_on  = [outscale_user.basic_user]
	}
	resource "outscale_user" "basic_user" {
		user_name = "%s"
		path = "/"
	}`, userName)
}

func testAccOutscaleUserUpdatedConfig(name string) string {
	return fmt.Sprintf(`
		resource "outscale_user" "update_user" {
			user_name = "%s"
			path = "/"
		}
	`, name)
}

func testAccUserWithPolicy(policyName, userName string) string {
	return fmt.Sprintf(`
		resource "outscale_policy" "policy-1" {
			policy_name = "%s"
		  	description = "testacc-user-terraform"
			document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
		  	path        = "/"
		}

		resource "outscale_user" "user_policy" {
			user_name = "%s"
			path = "/"
			policy {
	        	policy_orn = outscale_policy.policy-1.orn
			}
		}
	`, policyName, userName)
}

func testAccUserWithPolicyVersionUpper(policyName, userName string) string {
	return fmt.Sprintf(`
		resource "outscale_policy" "policy-2" {
			policy_name = "%s"
			description = "testacc-user-terraform"
			document    = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"*\"], \"Resource\": [\"*\"]} ]}"
			path        = "/"
		}

		resource "outscale_policy_version" "policy-2-v2" {
			policy_orn = outscale_policy.policy-2.orn
			document   = "{\"Statement\": [ {\"Effect\": \"Allow\", \"Action\": [\"eim:*\"], \"Resource\": [\"*\"]} ]}"
		}

		resource "outscale_user" "user_policy_version" {
			user_name = "%s"
			path = "/"
			policy {
				policy_orn = outscale_policy.policy-2.orn
				default_version_id = "V2"
			}
		}
	`, policyName, userName)
}
