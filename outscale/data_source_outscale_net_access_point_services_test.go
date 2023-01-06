package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutscaleOAPINetAccessPointServices_basic(t *testing.T) {
	t.Parallel()
	serviceName := fmt.Sprintf("com.outscale.%s.api", os.Getenv("OUTSCALE_REGION"))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPINetAccessPointServicesConfig(serviceName, serviceName),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPINetAccessPointServicesConfig(sName, sName2 string) string {
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
                        service_name    = "%s"
                        tags { 
                              key       = "name" 
                              value     = "terraform-Net-Access-Point" 
                        }

                }

               data "outscale_net_access_point_services" "all-services" {
                        filter {
                               name     = "service_names"
                               values   = [ "%s"]
                        }
               }

               data "outscale_net_access_point_services" "all-services2" { }

	`, sName, sName2)
}
