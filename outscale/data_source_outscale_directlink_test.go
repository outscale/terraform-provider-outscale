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

func TestAccOutscaleDSDL_basic(t *testing.T) {
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
				Config: testAccDLDSConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleDLDataSourceID("data.outscale_directlink.test"),
					resource.TestCheckResourceAttrSet("data.outscale_directlink.test", "connections.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleDLDataSourceID(n string) resource.TestCheckFunc {
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

func testAccDLDSConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`

		data "outscale_sites" "test" {}

		resource "outscale_directlink" "hoge" {
  			bandwidth = "1Gbps"
    		connection_name = "test-directlink-%d"
    		location = "${data.outscale_sites.test.locations.0.location_code}"
		}

		data "outscale_directlink" "test" {
			connection_id = "${outscale_directlink.hoge.id}"
		}
	`, rInt)
}
