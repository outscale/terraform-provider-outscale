package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceOutscaleOAPISecurityGroup_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_security_group.by_id"),
					//testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_security_group.by_filter"),
				),
			},
		},
	})
}

func TestAccDataSourceOutscaleOAPISecurityGroupPublic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISecurityGroupPublicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleOAPISecurityGroupCheck("data.outscale_security_group.by_filter_public"),
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

		SGRs, ok := s.RootModule().Resources["outscale_security_group.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_security_group.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["security_group_id"] != SGRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"security_group_id is %s; want %s",
				attr["security_group_id"],
				SGRs.Primary.Attributes["id"],
			)
		}

		//TODO: validate tags
		// if attr["tag.Name"] != "tf-acctest" {
		// 	return fmt.Errorf("bad Name tag %s", attr["tag.Name"])
		// }

		return nil
	}
}

func testAccDataSourceOutscaleOAPISecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "vpc" {
				ip_range = "10.0.0.0/16"
		}

		resource "outscale_security_group" "test" {
			net_id = "${outscale_net.vpc.id}"
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-%d"
			#tag = {
			#	Name = "tf-acctest"
			#	Seed = "%d"
			#}
		}

		data "outscale_security_group" "by_id" {
			security_group_id = "${outscale_security_group.test.id}"
		}

		#data "outscale_security_group" "by_filter" {
		#	filter {
		#		name = "security_group_names"
		#		values = ["${outscale_security_group.test.security_group_name}"]
		#	}
		#}`, rInt, rInt)
}

func testAccDataSourceOutscaleOAPISecurityGroupPublicConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "test" {
			description = "Used in the terraform acceptance tests"
			security_group_name = "test-%d"
			tag = {
				Name = "tf-acctest"
				Seed = "%d"
			}
		}

		data "outscale_security_group" "by_filter_public" {
			filter {
		    name = "security_group_names"

				// name = "group_name"
				values = ["${outscale_security_group.test.security_group_name}"]
			}
		}`, rInt, rInt)
}
