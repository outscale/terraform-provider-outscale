package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccKeypairsDataSource_Instance(t *testing.T) {
	t.Parallel()
	keyPairName := fmt.Sprintf("testacc-keypair-%d", acctest.RandIntRange(0, 400))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckKeypairsDataSourceConfig(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeypairsDataSourceID("data.outscale_keypairs.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypairs.nat_ami", "keypairs.0.keypair_name", keyPairName),
				),
			},
		},
	})
}

func testAccCheckKeypairsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Key Pair data source ID not set")
		}
		return nil
	}
}

func testAccCheckKeypairsDataSourceConfig(keyPairName string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
			keypair_name = "%s"
		}
		
		data "outscale_keypairs" "nat_ami" {
			#keypair_name = ["${outscale_keypair.a_key_pair.id}"]
		
			filter {
				name   = "keypair_names"
				values = ["${outscale_keypair.a_key_pair.keypair_name}"]
			}
		}
	`, keyPairName)
}
