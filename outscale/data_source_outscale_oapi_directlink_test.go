package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIDSDL_basic(t *testing.T) {
	t.Skip()

	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIDLDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDLDataSourceID("data.outscale_directlink.test"),
					resource.TestCheckResourceAttrSet("data.outscale_directlink.test", "directlinks.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIDLDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Directlink data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Directlink data source ID not set")
		}
		return nil
	}
}

func testAccOAPIDLDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		data "outscale_sites" "test" {}

		resource "outscale_directlink" "hoge" {
  			bandwidth = "1Gbps"
    		directlink_name = "test-directlink-%d"
    		site = "${data.outscale_sites.test.sites.0.code}"
		}

		data "outscale_directlink" "test" {
			directlink_id = "${outscale_directlink.hoge.id}"
		}
	`, rInt)
}
