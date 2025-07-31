package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_Ephemeral_keypair_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:   testAccEphemeralKeypairBasicConfig,
				PlanOnly: true,
			},
		},
	})
}

const testAccEphemeralKeypairBasicConfig = `
	ephemeral "outscale_keypair" "basic_ephemeral_keypair" {
		keypair_name = "ephemeral-keypair"
	}
`
