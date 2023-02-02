package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_LoadBalancer_DataSource(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")
	zone := fmt.Sprintf("%sa", region)
	dataSourceName := "data.outscale_load_balancer.test"
	dataSourcesName := "data.outscale_load_balancers.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAcc_LoadBalancer_DataSource_Config(zone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "subregion_names.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "subregion_names.0", zone),

					resource.TestCheckResourceAttr(dataSourcesName, "load_balancers.#", "1"),
				)},
		},
	})
}

func testAcc_LoadBalancer_DataSource_Config(zone string) string {
	return fmt.Sprintf(`
	resource "outscale_load_balancer" "dataLb" {
		subregion_names    = ["%s"]
		load_balancer_name = "data-terraform-elb"

		listeners {
			backend_port           = 8000
			backend_protocol       = "HTTP"
			load_balancer_port     = 80
			load_balancer_protocol = "HTTP"
		}

		tags {
			key   = "name"
			value = "baz"
		}
	}

	data "outscale_load_balancer" "test" {
		filter {
			name   = "load_balancer_names"
			values = [outscale_load_balancer.dataLb.id]
		}
	}

	data "outscale_load_balancers" "test" {
		filter {
			name   = "load_balancer_names"
			values = [outscale_load_balancer.dataLb.id]
		}
	}
`, zone)
}
