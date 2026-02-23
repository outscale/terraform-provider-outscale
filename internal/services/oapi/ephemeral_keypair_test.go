package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_Ephemeral_keypair_basic(t *testing.T) {
	keypairName := acctest.RandomWithPrefix("testacc-ephemeral-keypair")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

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
