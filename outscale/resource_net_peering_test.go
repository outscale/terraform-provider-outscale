package outscale

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_PeeringConnection_basic(t *testing.T) {
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

func TestAccNet_PeeringConnection_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig(utils.GetAccepterOwnerId()),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleLinkPeeeringConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleLinkPeeeringConnectionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
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
