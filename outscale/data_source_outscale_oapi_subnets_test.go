package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPISubnets(t *testing.T) {
	rInt := acctest.RandIntRange(0, 256)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_subnets", "subnet_set.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetsConfig(rInt int) string {
	return fmt.Sprintf(`
		
		resource "outscale_subnet" "test" {
		  vpc_id            = "vpc-e9d09d63"
		  cidr_block        = "172.%d.123.0/24"
		  availability_zone = "eu-west-2a"
		}
	
		data "outscale_subnets" "by_filter" {
		  filter {
		    name = "subnet-id"
		    values = ["${outscale_subnet.test.id}"]
		  }
		}
		`, rInt)
}
