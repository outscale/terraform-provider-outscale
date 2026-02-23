package oapi_test

import (
	"fmt"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_Peeringconnection_Basic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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

func TestAccNet_Peeringconnection_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig(oapihelpers.GetAccepterOwnerId()),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccNet_Peeringconnection_plan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:             testAccOAPIVpcPeeringConfig(oapihelpers.GetAccepterOwnerId()),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccNet_Peeringconnection_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.1.1",
			testAccOAPIVpcPeeringConfig2(),
			testAccOAPIVpcPeeringConfig(oapihelpers.GetAccepterOwnerId()),
		),
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
