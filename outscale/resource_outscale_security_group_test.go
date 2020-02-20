package outscale

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPISecurityGroup(t *testing.T) {
	var group oscgo.SecurityGroup
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPISGRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISecurityGroupRuleExists("outscale_security_group.web", &group),
					resource.TestCheckResourceAttr(
						"outscale_security_group.web", "security_group_name", fmt.Sprintf("terraform_test_%d", rInt)),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISGRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_security_group" {
			continue
		}

		// Retrieve our group
		req := oscgo.ReadSecurityGroupsRequest{
			Filters: &oscgo.FiltersSecurityGroup{
				SecurityGroupIds: &[]string{rs.Primary.ID},
			},
		}

		var resp oscgo.ReadSecurityGroupsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(context.Background(), &oscgo.ReadSecurityGroupsOpts{ReadSecurityGroupsRequest: optional.NewInterface(req)})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		var errString string

		if err != nil {
			if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
				return err
			}
			//fmt.Printf("\n\nError on SGStateRefresh: %s", err)
			errString = err.Error()

			return fmt.Errorf("Error on SGStateRefresh: %s", errString)
		}
		if err == nil {
			if len(resp.GetSecurityGroups()) > 0 && resp.GetSecurityGroups()[0].GetSecurityGroupId() == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if resp.GetSecurityGroups() == nil {
			return nil
		}

		return err
	}

	return nil
}

func testAccCheckOutscaleOAPISecurityGroupRuleExists(n string, group *oscgo.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		req := oscgo.ReadSecurityGroupsRequest{
			Filters: &oscgo.FiltersSecurityGroup{
				SecurityGroupIds: &[]string{rs.Primary.ID},
			},
		}

		var resp oscgo.ReadSecurityGroupsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(context.Background(), &oscgo.ReadSecurityGroupsOpts{ReadSecurityGroupsRequest: optional.NewInterface(req)})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		var errString string

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidSecurityGroupID.NotFound") ||
				strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
				err = nil
			} else {
				//fmt.Printf("\n\nError on SGStateRefresh: %s", err)
				errString = err.Error()
			}

			return fmt.Errorf("Error on SGStateRefresh: %s", errString)
		}

		if len(resp.GetSecurityGroups()) > 0 && resp.GetSecurityGroups()[0].GetSecurityGroupId() == rs.Primary.ID {
			*group = resp.GetSecurityGroups()[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccOutscaleOAPISecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"
		}
		
		resource "outscale_security_group" "web" {
			security_group_name = "terraform_test_%d"
			description         = "Used in the terraform acceptance tests"
		
			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		
			net_id = "${outscale_net.net.id}"
		}	
	`, rInt)
}
