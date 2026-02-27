package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_KeypairsDataSource_Instance(t *testing.T) {
	keyPairName := fmt.Sprintf("testacc-keypair-%d", utils.RandIntRange(0, 400))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
			return fmt.Errorf("can't find keypair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("key pair data source id not set")
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
