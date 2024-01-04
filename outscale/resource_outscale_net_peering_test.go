package outscale

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_PeeringConnection_basic(t *testing.T) {
	var connection oscgo.NetPeering

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_net_peering.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinPeeringConnectionExists(
						"outscale_net_peering.foo",
						&connection),
				),
			},
		},
	})
}

func TestAccNet_PeeringConnection_importBasic(t *testing.T) {
	resourceName := "outscale_net_peering.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPILinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig,
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleOAPILinkPeeeringConnectionImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleOAPILinkPeeeringConnectionImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func TestAccNet_PeeringConnection_plan(t *testing.T) {
	var connection oscgo.NetPeering

	// reach out and DELETE the VPC Peering connection outside of Terraform
	testDestroy := func(*terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		log.Printf("[DEBUG] Test deleting the Net Peering.")
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.NetPeeringApi.DeleteNetPeering(context.Background()).DeleteNetPeeringRequest(oscgo.DeleteNetPeeringRequest{
				NetPeeringId: connection.GetNetPeeringId(),
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPILinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIVpcPeeringConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinPeeringConnectionExists(
						"outscale_net_peering.foo",
						&connection),
					testDestroy,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckOutscaleOAPILinPeeringConnectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_net_peering" {
			continue
		}

		var resp oscgo.ReadNetPeeringsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(oscgo.ReadNetPeeringsRequest{
				Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{rs.Primary.ID}},
			}).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		var errString string

		if err != nil {
			errString = err.Error()
			return fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		pc := &oscgo.NetPeering{}
		for _, c := range resp.GetNetPeerings() {
			if rs.Primary.ID == c.GetNetPeeringId() {
				pc = &c
			}
		}

		if pc == nil {
			// not found
			return nil
		}

		if pc.State.GetName() != "" {
			if pc.State.GetName() == "deleted" {
				return nil
			}
			return fmt.Errorf("Found the Net Peering in an unexpected state: %v", pc)
		}

		// return error here; we've found the vpc_peering object we want, however
		// it's not in an expected state
		return fmt.Errorf("Fall through error for testAccCheckOutscaleOAPILinPeeringConnectionDestroy")
	}
	return nil
}

func testAccCheckOutscaleOAPILinPeeringConnectionExists(n string, connection *oscgo.NetPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Net Peering ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		var resp oscgo.ReadNetPeeringsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(oscgo.ReadNetPeeringsRequest{
				Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{rs.Primary.ID}},
			}).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		var errString string

		if err != nil {
			errString = err.Error()
			return fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		if len(resp.GetNetPeerings()) == 0 {
			return fmt.Errorf("Net Peering could not be found")
		}

		*connection = resp.GetNetPeerings()[0]

		return nil
	}
}

const testAccOAPIVpcPeeringConfig = `
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
