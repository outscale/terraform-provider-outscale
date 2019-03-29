package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPISecurityGroups_vpc(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISecurityGroupConfigVPC(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_security_groups.by_id", "security_groups.#", "3"),
					resource.TestCheckResourceAttr(
						"data.outscale_security_groups.by_filter", "security_groups.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISecurityGroupConfigVPC(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_security_group" "test" {
		net_id = "vpc-e9d09d63"
		description = "Used in the terraform acceptance tests"
		security_group_name = "test-1--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	resource "outscale_security_group" "test2" {
		net_id = "vpc-e9d09d63"
		description = "Used in the terraform acceptance tests"
		security_group_name = "test-2--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	resource "outscale_security_group" "test3" {
		net_id = "vpc-e9d09d63"
		description = "Used in the terraform acceptance tests"
		security_group_name = "test-3--%d"
		tag = {
			Name = "tf-acctest"
			Seed = "%d"
		}
	}
	data "outscale_security_groups" "by_id" {
		security_group_id = ["${outscale_security_group.test.id}", "${outscale_security_group.test2.id}", "${outscale_security_group.test3.id}"]
	}
	data "outscale_security_groups" "by_filter" {
		filter {
			name = "group-name"
			values = ["${outscale_security_group.test.security_group_name}"]
		}
	}
	`, rInt, rInt, rInt, rInt, rInt, rInt)
}
