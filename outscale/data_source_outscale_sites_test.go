package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSSites_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPICustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPISitesDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISitesDataSourceID("data.outscale_sites.test"),
					resource.TestCheckResourceAttrSet("data.outscale_sites.test", "sites.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISitesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Sites data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Sites data source ID not set")
		}
		return nil
	}
}

func testAccOAPISitesDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_sites" "test" {}
	`)
}
