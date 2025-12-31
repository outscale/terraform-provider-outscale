package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_AccessPoint_Basic(t *testing.T) {
	serviceName := fmt.Sprintf("com.outscale.%s.api", utils.GetRegion())
	resourceName := "outscale_net_access_point.net_access_point_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testacc.PreCheck(t)
		},
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNetAccessPointConfig(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "available"),
				),
			},
		},
	})
}

func TestAccNet_AccessPoint_import(t *testing.T) {
	resourceName := "outscale_net_access_point.net_access_point_1"
	serviceName := fmt.Sprintf("com.outscale.%s.api", utils.GetRegion())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNetAccessPointConfig(serviceName),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_AccessPoint_Migration(t *testing.T) {
	serviceName := fmt.Sprintf("com.outscale.%s.api", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.1.3", testAccOutscaleNetAccessPointConfig(serviceName)),
	})
}

func testAccOutscaleNetAccessPointConfig(sName string) string {
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
	`, sName)
}
