package outscale

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_PeeringConnection_Basic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "accepter_net_id"),
					resource.TestCheckResourceAttr(resourceName, "state.0.name", "pending-acceptance"),
				),
			},
		},
	})
}

func TestAccNet_PeeringConnection_Basic_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.1", testAccOAPIVpcPeeringConfig2()),
	})
}

func TestAccNet_PeeringConnection_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig(utils.GetAccepterOwnerId()),
			},
			testutils.ImportStepFW(resourceName, testutils.DefaultIgnores()...),
		},
	})
}

func TestAccNet_PeeringConnection_importBasic_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps:    FrameworkMigrationTestSteps("1.1.1", testAccOAPIVpcPeeringConfig(utils.GetAccepterOwnerId())),
	})
}

func TestAccNet_PeeringConnection_plan(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config:             testAccOAPIVpcPeeringConfig(utils.GetAccepterOwnerId()),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func testAccOAPIVpcPeeringConfig(accountid string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "foo" {
		ip_range = "10.0.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-rs-foo"
		}
	}

	resource "outscale_net" "bar" {
		ip_range = "10.1.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-acceptation-rs-bar"
		}
	}

	resource "outscale_net_peering" "foo" {
		source_net_id   = outscale_net.foo.id
		accepter_net_id = outscale_net.bar.id
		accepter_owner_id = "%s"
	}
`, accountid)
}

func testAccOAPIVpcPeeringConfig2() string {
	return `
	resource "outscale_net" "foo" {
		ip_range = "10.0.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-rs-foo"
		}
	}

	resource "outscale_net" "bar" {
		ip_range = "10.1.0.0/16"

		tags {
			key   = "Name"
			value = "testacc-net-peering-acceptation-rs-bar"
		}
	}

	resource "outscale_net_peering" "foo" {
		source_net_id   = outscale_net.foo.id
		accepter_net_id = outscale_net.bar.id
	}
`
}
