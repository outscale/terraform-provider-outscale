package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testAccCheckOutscaleSecurityGroupExists(n string, group *oscgo.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		fids := []string{rs.Primary.ID}
		filter := oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &fids,
		}

		req := &oscgo.ReadSecurityGroupsRequest{
			Filters: &filter,
		}
		var err error
		var resp oscgo.ReadSecurityGroupsResponse
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(
				context.Background(),
				&oscgo.ReadSecurityGroupsOpts{
					ReadSecurityGroupsRequest: optional.NewInterface(req)})

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

		if resp.SecurityGroups != nil && len(*resp.SecurityGroups) > 0 &&
			*(*resp.SecurityGroups)[0].SecurityGroupId == rs.Primary.ID {
			*group = (*resp.SecurityGroups)[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccCheckOutscaleSecurityGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_firewall_rules_set" {
			continue
		}

		// Retrieve our group
		fids := []string{rs.Primary.ID}
		filter := oscgo.FiltersSecurityGroup{
			SecurityGroupIds: &fids,
		}

		req := &oscgo.ReadSecurityGroupsRequest{
			Filters: &filter,
		}

		var err error
		var resp oscgo.ReadSecurityGroupsResponse
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			resp, _, err = conn.SecurityGroupApi.ReadSecurityGroups(
				context.Background(),
				&oscgo.ReadSecurityGroupsOpts{
					ReadSecurityGroupsRequest: optional.NewInterface(req)})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			if resp.SecurityGroups != nil && len(*resp.SecurityGroups) > 0 &&
				*(*resp.SecurityGroups)[0].SecurityGroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
			return nil
		}
	}

	return nil
}

const testAccOutscaleSecurityGroupConfigClassic = `
resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_1"
  group_description = "Used in the terraform acceptance tests"
}
`
