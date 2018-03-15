package outscale

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleKeyPairDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeyPairDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					//testAccCheckOutscaleKeyPairExists("outscale_keypair.a_key_pair", &conf),
					resource.TestCheckResourceAttr("data.outscale_keypair", "key_name", "tf-acc-key-pair2"),
				),
			},
		},
	})
}

const testAccCheckOutscaleKeyPairDataSourceConfigBasic = `
resource "outscale_keypair" "a_key_pair" {
	key_name   = "tf-acc-key-pair2"
}

data "outscale_keypair" "outscale_keypair" {
    filter {
	name = "key-name"
	values = ["${outscale_keypair.a_key_pair.key_name}"]
    }
}
`
