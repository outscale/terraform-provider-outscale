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
					resource.TestCheckResourceAttr("data.outscale_subnets.by_filter", "subnets.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPISubnetsConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "outscale_net" "net" {
		ip_range = "172.%[1]d.123.0/24"
	}

	resource "outscale_subnet" "subnet" {
		ip_range = "172.%[1]d.123.0/24"
		subregion_name = "eu-west-2a"
		net_id = "${outscale_net.net.id}"
	
		tags = {
			key = "name"
			value = "terraform-subnet"
		}
		}
	
		data "outscale_subnets" "by_filter" {
		  filter {
		    name = "subnet_ids"
		    values = ["${outscale_subnet.subnet.id}"]
		  }
		}
		`, rInt)
}
