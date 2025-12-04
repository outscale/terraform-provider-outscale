package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccOthers_KeypairsDataSource_Instance(t *testing.T) {
	t.Parallel()
	keyPairName := fmt.Sprintf("testacc-keypair-%d", utils.RandIntRange(0, 400))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeypairsDataSourceConfig(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeypairsDataSourceID("data.outscale_keypairs.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypairs.nat_ami", "keypairs.0.keypair_name", keyPairName),
				),
			},
		},
	})
}

func testAccCheckOutscaleKeypairsDataSourceID(n string) resource.TestCheckFunc {
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

func testAccCheckOutscaleKeypairsDataSourceConfig(keyPairName string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
			keypair_name = "%s"
		}

		data "outscale_keypairs" "nat_ami" {
			filter {
				name   = "keypair_names"
				values = [outscale_keypair.a_key_pair.keypair_name]
			}
		}
	`, keyPairName)
}
