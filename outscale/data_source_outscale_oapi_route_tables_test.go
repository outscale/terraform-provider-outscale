package outscale

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIRouteTables_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIRouteTablesGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_route_tables.by_filter", "route_table_set.#", "1"),
					resource.TestCheckResourceAttr("data.outscale_route_tables.by_id", "route_table_set.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceOutscaleOAPIRouteTablesGroupConfig = `
resource "outscale_lin" "test" {
  ip_range = "172.16.0.0/16"

  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_subnet" "test" {
  ip_range = "172.16.0.0/24"
  lin_id     = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_route_table" "test" {
  lin_id = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-routetable-data-source"
  }
}

data "outscale_route_tables" "by_filter" {
  filter {
    name = "route-table-id"
    values = ["${outscale_route_table.test.id}"]
  }
}

data "outscale_route_tables" "by_id" {
  route_table_id = ["${outscale_route_table.test.id}"]
}
`
