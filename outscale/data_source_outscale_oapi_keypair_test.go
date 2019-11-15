package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIKeypairDataSource_Instance(t *testing.T) {
	t.Skip()

	keyPairName := fmt.Sprintf("acc-test-keypair-%d", acctest.RandIntRange(0, 400))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIKeypairDataSourceConfig(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIKeypairDataSourceID("data.outscale_keypair.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypair.nat_ami", "keypair_name", keyPairName),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIKeypairDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find key pair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Key Pair data source ID not set")
		}
		return nil
	}
}

func testAccCheckOutscaleOAPIKeypairDataSourceConfig(keypairName string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
			keypair_name = "%s"
		}
		
		data "outscale_keypair" "nat_ami" {
			#keypair_name = "${outscale_keypair.a_key_pair.id}"
			filter {
				name   = "keypair_names"
				values = ["${outscale_keypair.a_key_pair.keypair_name}"]
			}
		}
	`, keypairName)
}
