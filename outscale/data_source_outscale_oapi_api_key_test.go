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

func TestAccOutscaleDSOAPIAPIKey_basic(t *testing.T) {
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIAPIKeyDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIAKsDataSourceID("data.outscale_api_key.test"),
					resource.TestCheckResourceAttr("data.outscale_api_key.test", "access_key.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIAKsDataSourceID(n string) resource.TestCheckFunc {
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

func testAccOAPIAPIKeyDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_api_key" "test" {}
	`)
}
