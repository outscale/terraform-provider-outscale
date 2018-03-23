package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleSubnets(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 256)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleSubnetsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_subnets.by_filter", "subnet_set.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleSubnetsConfig(rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_lin" "vpc" {
			cidr_block = "10.0.0.0/16"
		}

		resource "outscale_subnet" "test" {
		  vpc_id            = "${outscale_lin.vpc.id}"
		  cidr_block        = "10.0.%d.0/24"
		  availability_zone = "eu-west-2a"
		}
	
		data "outscale_subnets" "by_filter" {}
		`, rInt)
}
