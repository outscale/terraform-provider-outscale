package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleDSOAPIAccountConsumption_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIAccountConsumptionDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIAccountConsumptionDataSourceID("data.outscale_account_consumption.test"),
					resource.TestCheckResourceAttrSet("data.outscale_account_consumption.test", "consumption_entries.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIAccountConsumptionDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Acctount Consumption data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Account Consumption data source ID not set")
		}
		return nil
	}
}

func testAccOAPIAccountConsumptionDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_account_consumption" "test" {
			from_date = "2018-02-01"
			to_date = "2018-07-01"
		}
	`)
}
