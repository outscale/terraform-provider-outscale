package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_Ephemeral_keypair_basic(t *testing.T) {
	keypairName := acctest.RandomWithPrefix("testacc-ephemeral-keypair")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:   testAccEphemeralKeypairBasicConfig(keypairName),
				PlanOnly: true,
			},
		},
	})
}

func testAccEphemeralKeypairBasicConfig(keypairName string) string {
	return fmt.Sprintf(`
	ephemeral "outscale_keypair" "basic_ephemeral_keypair" {
		keypair_name = "%s"
	}
`, keypairName)
}
