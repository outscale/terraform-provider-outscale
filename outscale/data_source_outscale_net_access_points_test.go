package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_AccessPointsDataSource_basic(t *testing.T) {
	t.Parallel()
	serviceName := fmt.Sprintf("com.outscale.%s.api", utils.GetRegion())
	resourceName := "data.outscale_net_access_points.data_net_access_points"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleNetAccessPointsConfig(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_access_points.0.route_table_ids.#"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleNetAccessPointsConfig(sName string) string {
	return fmt.Sprintf(`
                resource "outscale_net" "outscale_net" {
                        ip_range = "10.0.0.0/16"
                }

                resource "outscale_route_table" "route_table-1" {
                        net_id = outscale_net.outscale_net.net_id
                }

                resource "outscale_net_access_point" "net_access_point_1" {
                        net_id          = outscale_net.outscale_net.net_id
                        route_table_ids = [outscale_route_table.route_table-1.route_table_id]
                        service_name    = "%[1]s"
                        tags {
                              key       = "name"
                              value     = "terraform-Net-Access-Point"
                        }

                }

               data "outscale_net_access_points" "data_net_access_points" {
                        filter {
                               name     = "net_access_point_ids"
                               values   = [outscale_net_access_point.net_access_point_1.net_access_point_id]
                        }
                        filter {
                               name     = "net_ids"
                               values   = [outscale_net.outscale_net.net_id]
                        }
                        filter {
                               name     = "service_names"
                               values   = [ "%[1]s"]
                        }
                        filter {
                               name     = "states"
                               values   = ["available"]
                        }
                        filter {
                               name     = "tag_keys"
                               values   = ["name"]
                        }
                        filter {
                               name     = "tag_values"
                               values   = ["terraform-Net-Access-Point"]
                        }
               }

	`, sName)
}
