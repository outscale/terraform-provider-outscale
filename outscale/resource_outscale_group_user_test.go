package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleGroupUser_basic(t *testing.T) {
	var group eim.GetGroupOutput

	rInt := acctest.RandInt()
	configBase := fmt.Sprintf(testAccOutscaleGroupUserConfig, rInt, rInt)
	testUser := fmt.Sprintf("test-user-%d", rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleGroupUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleGroupUserExists("outscale_group_user.team", &group),
					testAccCheckOutscaleGroupUserAttributes(&group, []string{testUser}),
				),
			},
		},
	})
}

func testAccCheckOutscaleGroupUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_group_user" {
			continue
		}

		group := rs.Primary.Attributes["group_name"]

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.API.GetGroup(&eim.GetGroupInput{
				GroupName: aws.String(group),
			})

			if err != nil {
				if strings.Contains(err.Error(), "Throttling:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
				return nil
			}
			return err
		}

		return fmt.Errorf("still exists")
	}

	return nil
}

func testAccCheckOutscaleGroupUserExists(n string, g *eim.GetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User name is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).EIM
		gn := rs.Primary.Attributes["group_name"]

		var err error
		var resp *eim.GetGroupOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.GetGroup(&eim.GetGroupInput{
				GroupName: aws.String(gn),
			})
			if err != nil {
				if strings.Contains(err.Error(), "Throttling:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error: Group (%s) not found", gn)
		}

		*g = *resp

		return nil
	}
}

func testAccCheckOutscaleGroupUserAttributes(group *eim.GetGroupOutput, users []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !strings.Contains(*group.GetGroupResult.Group.GroupName, "test-group") {
			return fmt.Errorf("Bad group membership: expected %s, got %s", "test-group", *group.GetGroupResult.Group.GroupName)
		}

		return nil
	}
}

const testAccOutscaleGroupUserConfig = `
resource "outscale_group" "group" {
	group_name = "test-group-%d"
	path = "/"
}

resource "outscale_user" "user" {
	user_name = "test-user-%d"
	path = "/"
}

resource "outscale_group_user" "team" {
	user_name = "${outscale_user.user.user_name}"
	group_name = "${outscale_group.group.group_name}"
}
`
