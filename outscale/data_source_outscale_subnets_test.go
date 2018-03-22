package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleSubnets(t *testing.T) {
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
		
		resource "outscale_subnet" "test" {
		  vpc_id            = "vpc-e9d09d63"
		  cidr_block        = "10.0.%d.0/24"
		  availability_zone = "eu-west-2a"
		}
	
		data "outscale_subnets" "by_filter" {
		  #subnet_id = ["${outscale_subnet.test.id}"]
		}
		`, rInt)
}
