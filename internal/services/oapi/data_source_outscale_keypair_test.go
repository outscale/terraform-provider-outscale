package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_KeypairDataSource_Instance(t *testing.T) {
	keyPairName := fmt.Sprintf("acc-test-keypair-%d", utils.RandIntRange(0, 400))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeypairDataSourceConfig(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeypairDataSourceID("data.outscale_keypair.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypair.nat_ami", "keypair_name", keyPairName),
				),
			},
		},
	})
}

func testAccCheckOutscaleKeypairDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find key pair data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("key pair data source id not set")
		}
		return nil
	}
}

func testAccCheckOutscaleKeypairDataSourceConfig(keypairName string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
			keypair_name = "%s"
		}

		data "outscale_keypair" "nat_ami" {
			filter {
				name   = "keypair_names"
				values = [outscale_keypair.a_key_pair.keypair_name]
			}
		}
	`, keypairName)
}
