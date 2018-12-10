package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISecurityGroup(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if isOapi == false {
		t.Skip()
	}
	var group oapi.SecurityGroup
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
					testAccCheckState("outscale_security_group.web"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISGRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_security_group" {
			continue
		}

		// Retrieve our group
		req := &oapi.ReadSecurityGroupsRequest{
			Filters: oapi.FiltersSecurityGroup{
				SecurityGroupIds: []string{rs.Primary.ID},
			},
		}

		var resp *oapi.POST_ReadSecurityGroupsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadSecurityGroups(*req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				if strings.Contains(err.Error(), "InvalidGroup.NotFound") {
					return err
				} else {
					//fmt.Printf("\n\nError on SGStateRefresh: %s", err)
					errString = err.Error()
				}

			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("Error on SGStateRefresh: %s", errString)
		}

		result := resp.OK

		if err == nil {
			if len(result.SecurityGroups) > 0 && result.SecurityGroups[0].SecurityGroupId == rs.Primary.ID {
				return fmt.Errorf("Security Group (%s) still exists", rs.Primary.ID)
			}

			return nil
		}

		if result == nil {
			return nil
		}

		return err
	}

	return nil
}

func testAccCheckOutscaleOAPISecurityGroupRuleExists(n string, group *oapi.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		req := &oapi.ReadSecurityGroupsRequest{
			Filters: oapi.FiltersSecurityGroup{
				SecurityGroupIds: []string{rs.Primary.ID},
			},
		}

		var resp *oapi.POST_ReadSecurityGroupsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadSecurityGroups(*req)

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

		if err != nil || resp.OK == nil {
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidSecurityGroupID.NotFound") ||
					strings.Contains(fmt.Sprint(err), "InvalidGroup.NotFound") {
					resp = nil
					err = nil
				} else {
					//fmt.Printf("\n\nError on SGStateRefresh: %s", err)
					errString = err.Error()
				}

			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("Error on SGStateRefresh: %s", errString)
		}

		result := resp.OK

		if len(result.SecurityGroups) > 0 && result.SecurityGroups[0].SecurityGroupId == rs.Primary.ID {
			*group = result.SecurityGroups[0]
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

func testAccOutscaleOAPISecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_security_group" "web" {
		security_group_name = "terraform_test_%d"
		description = "Used in the terraform acceptance tests"
		tag = {
			Name = "tf-acc-test"
		}
		net_id = "vpc-e9d09d63"
	}`, rInt)
}
