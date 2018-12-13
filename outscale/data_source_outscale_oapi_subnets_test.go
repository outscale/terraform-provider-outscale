package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPISubnets(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 256)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPISubnetsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_subnets", "subnets.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetsConfig(rInt int) string {
	return fmt.Sprintf(`
		
		resource "outscale_subnet" "test" {
		  net_id            = "vpc-e9d09d63"
		  ip_range        = "172.%d.123.0/24"
		  subregion_name = "eu-west-2a"
		}
	
		data "outscale_subnets" "by_filter" {
		  filter {
		    name = "subnet_ids"
		    values = ["${outscale_subnet.test.id}"]
		  }
		}
		`, rInt)
}
