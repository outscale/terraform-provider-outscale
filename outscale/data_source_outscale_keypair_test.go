package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAcc_Keypair_DataSource(t *testing.T) {
	t.Parallel()
	keyPairName := fmt.Sprintf("acc-test-keypair-%d", utils.RandIntRange(0, 400))
	dataSourceName := "data.outscale_keypair.keypair"
	dataSourcesName := "data.outscale_keypairs.keypairs"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAcc_Keypair_DataSource_Config(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "keypair_name", keyPairName),
					resource.TestCheckResourceAttr(dataSourcesName, "keypairs.#", "1"),
				),
			},
		},
	})
}

func testAcc_Keypair_DataSource_Config(keypairName string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "a_key_pair" {
			keypair_name = "%s"
		}
		
		data "outscale_keypair" "keypair" {
			filter {
				name   = "keypair_names"
				values = [outscale_keypair.a_key_pair.keypair_name]
			}
		}
		
		data "outscale_keypairs" "keypairs" {
			filter {
				name   = "keypair_names"
				values = [outscale_keypair.a_key_pair.keypair_name]
			}
		}
	`, keypairName)
}
