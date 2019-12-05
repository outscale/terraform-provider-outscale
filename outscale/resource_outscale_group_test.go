package outscale

import (
	"errors"
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

func TestAccOutscaleOAPIGroup_basic(t *testing.T) {
	t.Skip()

	var conf eim.GetGroupOutput
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIGroupExists("outscale_group.group", &conf),
					testAccCheckOutscaleOAPIGroupAttributes(&conf, fmt.Sprintf("test-group-%d", rInt), "/"),
				),
			},
			{
				Config: testAccOutscaleOAPIGroupConfig2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIGroupExists("outscale_group.group2", &conf),
					testAccCheckOutscaleOAPIGroupAttributes(&conf, fmt.Sprintf("test-group-%d-2", rInt), "/funnypath/"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIGroupDestroy(s *terraform.State) error {
	eimconn := testAccProvider.Meta().(*OutscaleClient).EIM

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_group" {
			continue
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = eimconn.API.GetGroup(&eim.GetGroupInput{
				GroupName: aws.String(rs.Primary.ID),
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err == nil {
			return errors.New("still exist")
		}

		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			return nil
		}
	}

	return nil
}

func testAccCheckOutscaleOAPIGroupExists(n string, res *eim.GetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Group name is set")
		}

		eimconn := testAccProvider.Meta().(*OutscaleClient).EIM

		var err error
		var resp *eim.GetGroupOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = eimconn.API.GetGroup(&eim.GetGroupInput{
				GroupName: aws.String(rs.Primary.ID),
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		*res = *resp

		return nil
	}
}

func testAccCheckOutscaleOAPIGroupAttributes(group *eim.GetGroupOutput, name string, path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *group.GetGroupResult.Group.GroupName != name {
			return fmt.Errorf("Bad name: %s when %s was expected", *group.GetGroupResult.Group.GroupName, name)
		}

		if *group.GetGroupResult.Group.Path != path {
			return fmt.Errorf("Bad path: %s when %s was expected", *group.GetGroupResult.Group.Path, path)
		}

		return nil
	}
}

func testAccOutscaleOAPIGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_group" "group" {
			group_name = "test-group-%d"
			path       = "/"
		}
	`, rInt)
}

func testAccOutscaleOAPIGroupConfig2(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_group" "group2" {
			group_name = "test-group-%d-2"
			path       = "/funnypath/"
		}
	`, rInt)
}
