package outscale

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVpcs_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(16)
	ipRange := fmt.Sprintf("172.%d.0.0/16", rInt)
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleOAPIVpcsConfig(ipRange, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_nets.by_id", "net.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVpcsConfig(ipRange, tag string) string {
	return fmt.Sprintf(`

resource "outscale_net" "test" {
  ip_range = "%s"
	
	#Not supported yet
  #tag {
  #  Name = "%s"
  #}
}

data "outscale_nets" "by_id" {
  net_id = ["${outscale_net.test.id}"]
}`, ipRange, tag)
}
