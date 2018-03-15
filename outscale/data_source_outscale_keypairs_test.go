package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleKeyPairsDataSource_basic(t *testing.T) {
	var conf fcu.KeyPairInfo

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeyPairDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeyPairExists("outscale_keypair.a_key_pair", &conf),
					resource.TestCheckResourceAttr("data.outscale_keypairs", "key_set.#", "1"),
				),
			},
		},
	})
}

const testAccCheckOutscaleKeyPairsDataSourceConfigBasic = `
resource "outscale_keypair" "a_key_pair" {
	key_name   = "tf-acc-key-pair"
}

data "outscale_keypairs" "outscale_keypairs" {
    filter {
	name = "key-name"
	values = ["${outscale_keypair.a_key_pair.key_name}"]
    }
}
`
