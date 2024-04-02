package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithSecurityGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_security_group.by_id"),
					//testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_security_group.by_filter"),
				),
			},
		},
	})
}

func TestAccNet_WithSecurityGroupPublic_(t *testing.T) {
	t.Parallel()
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscaleSecurityGroupCheck("data.outscale_security_group.by_filter_public"),
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

func testAccDataSourceOutscaleSecurityGroupConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_net" "vpc" {
				ip_range = "10.0.0.0/16"
				tags {
					key = "Name"
					value = "testacc-sec-group-ds"
				}
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
		#		values = [outscale_security_group.test.security_group_name]
		#	}
		#}`, rInt, rInt)
}

func testAccDataSourceOutscaleSecurityGroupPublicConfig(rInt int) string {
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
				values = [outscale_security_group.test.security_group_name]
			}
		}`, rInt, rInt)
}
