package outscale

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleVpcs_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(16)
	cidr := fmt.Sprintf("172.%d.0.0/16", rInt)
	tag := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscaleVpcsConfig(cidr, tag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.outscale_lins.by_id", "vpc_set.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleVpcsConfig(cidr, tag string) string {
	return fmt.Sprintf(`

resource "outscale_lin" "test" {
  cidr_block = "%s"

  tag {
    Name = "%s"
  }
}

data "outscale_lins" "by_id" {
  vpc_id = ["${outscale_lin.test.id}"]
}`, cidr, tag)
}
