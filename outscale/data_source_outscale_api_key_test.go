package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleDSAPIKey_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAPIKeyDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAKsDataSourceID("data.outscale_api_key.test"),
					resource.TestCheckResourceAttr("data.outscale_api_key.test", "access_key.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleAKsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find API Key data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("API Key data source ID not set")
		}
		return nil
	}
}

func testAccAPIKeyDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_api_key" "test" {}
	`)
}
