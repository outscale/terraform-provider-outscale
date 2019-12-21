package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPILinPeeringConnection_basic(t *testing.T) {
	var connection oscgo.NetPeering

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_net_peering.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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

func TestAccOutscaleOAPILinPeeringConnection_plan(t *testing.T) {
	var connection oscgo.NetPeering

	// reach out and DELETE the VPC Peering connection outside of Terraform
	testDestroy := func(*terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		log.Printf("[DEBUG] Test deleting the Net Peering.")
		_, _, err := conn.NetPeeringApi.DeleteNetPeering(context.Background(), &oscgo.DeleteNetPeeringOpts{DeleteNetPeeringRequest: optional.NewInterface(oscgo.DeleteNetPeeringRequest{
			NetPeeringId: connection.GetNetPeeringId(),
		})})
		if err != nil {
			return err
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPILinPeeringConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
			resp, _, err = conn.NetPeeringApi.ReadNetPeerings(context.Background(), &oscgo.ReadNetPeeringsOpts{ReadNetPeeringsRequest: optional.NewInterface(oscgo.ReadNetPeeringsRequest{
				Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{rs.Primary.ID}},
			})})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
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
			resp, _, err = conn.NetPeeringApi.ReadNetPeerings(context.Background(), &oscgo.ReadNetPeeringsOpts{ReadNetPeeringsRequest: optional.NewInterface(oscgo.ReadNetPeeringsRequest{
				Filters: &oscgo.FiltersNetPeering{NetPeeringIds: &[]string{rs.Primary.ID}},
			})})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
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
			value = "TestAccOutscaleOAPILinPeeringConnection_basic"
		}
	}

	resource "outscale_net" "bar" {
		ip_range = "10.1.0.0/16"
	}

	resource "outscale_net_peering" "foo" {
		source_net_id   = "${outscale_net.foo.id}"
		accepter_net_id = "${outscale_net.bar.id}"
	}
`

//FIXME: check where is used.
// func testAccCheckOutscaleOAPILinPeeringConnectionOptions(n, block string, options *oscgo.NetPeeringOptionsDescription) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", n)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No Net Peering ID is set")
// 		}

// 		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

// 		var resp *oscgo.ReadNetPeeringsOutput
// 		var err error
// 		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			resp, err = conn.VM.ReadNetPeerings(
// 				&oscgo.ReadNetPeeringsRequest{
// 					NetPeeringIds: []*string{aws.String(rs.Primary.ID)},
// 				})

// 			if err != nil {
// 				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
// 					return resource.RetryableError(err)
// 				}
// 				return resource.NonRetryableError(err)
// 			}
// 			return nil
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		pc := resp.NetPeerings[0]

// 		o := pc.AccepterVpcInfo
// 		if block == "requester_vpc_info" {
// 			o = pc.RequesterVpcInfo
// 		}

// 		if !reflect.DeepEqual(o.PeeringOptions, options) {
// 			return fmt.Errorf("Expected the Net Peering Options to be %#v, got %#v",
// 				options, o.PeeringOptions)
// 		}

// 		return nil
// 	}
// }
