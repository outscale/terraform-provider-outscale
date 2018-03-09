package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPISecurityGroup(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_firewall_rule_set.by_id"),
					testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_firewall_rule_set.by_filter"),
				),
			},
		},
	})
}
func TestAccDataSourceOutscaleOAPISecurityGroupPublic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISecurityGroupPublicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_firewall_rule_set.by_filter_public"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISecurityGroupCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		SGRs, ok := s.RootModule().Resources["outscale_firewall_rules_set.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_firewall_rules_set.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["firewall_rules_set_id"] != SGRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"firewall_rules_set_id is %s; want %s",
				attr["firewall_rules_set_id"],
				SGRs.Primary.Attributes["id"],
			)
		}

		if attr["tag.Name"] != "tf-acctest" {
			return fmt.Errorf("bad Name tag %s", attr["tag.Name"])
		}

		return nil
	}
}

func testAccDataSourceOutscaleOAPISecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "test" {
		vpc_id = "vpc-e9d09d63"
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_firewall_rule_set" "by_id" {
		firewall_rules_set_id = "${outscale_firewall_rules_set.test.id}"
	}
	data "outscale_firewall_rule_set" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}`, rInt, rInt)
}

func testAccDataSourceOutscaleOAPISecurityGroupPublicConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "test" {
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_firewall_rule_set" "by_filter_public" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}`, rInt, rInt)
}
