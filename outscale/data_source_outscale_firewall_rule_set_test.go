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

func TestAccDataSourceOutscaleSecurityGroup(t *testing.T) {
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
				Config: testAccDataSourceOutscaleSecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_firewall_rules_set.by_id"),
					testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_firewall_rules_set.by_filter"),
				),
			},
		},
	})
}
func TestAccDataSourceOutscaleSecurityGroupPublic(t *testing.T) {
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
				Config: testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_firewall_rules_set.by_filter_public"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSecurityGroupCheck(name string) resource.TestCheckFunc {
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

		if attr["group_id"] != SGRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"group_id is %s; want %s",
				attr["group_id"],
				SGRs.Primary.Attributes["id"],
			)
		}

		if attr["tag_set.Name"] != "tf-acctest" {
			return fmt.Errorf("bad Name tag %s", attr["tag_set.Name"])
		}

		return nil
	}
}

func testAccDataSourceOutscaleSecurityGroupConfig(rInt int) string {
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
	data "outscale_firewall_rules_set" "by_id" {
		group_id = "${outscale_firewall_rules_set.test.id}"
	}
	data "outscale_firewall_rules_set" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}`, rInt, rInt)
}

func testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_firewall_rules_set" "test" {
		group_description = "Used in the terraform acceptance tests"
		group_name = "test-%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_firewall_rules_set" "by_filter_public" {
		filter {
			name = "group-name"
			values = ["${outscale_firewall_rules_set.test.group_name}"]
		}
	}`, rInt, rInt)
}
