package outscale

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPILinPeeringConnection_basic(t *testing.T) {
	var connection oapi.NetPeering

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
	var connection oapi.NetPeering

	// reach out and DELETE the VPC Peering connection outside of Terraform
	testDestroy := func(*terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		log.Printf("[DEBUG] Test deleting the Net Peering.")
		_, err := conn.POST_DeleteNetPeering(oapi.DeleteNetPeeringRequest{
			NetPeeringId: connection.NetPeeringId,
		})
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
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_net_peering" {
			continue
		}

		var resp *oapi.POST_ReadNetPeeringsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadNetPeerings(oapi.ReadNetPeeringsRequest{
				Filters: oapi.FiltersNetPeering{NetPeeringIds: []string{rs.Primary.ID}},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
			}
			return fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		result := resp.OK

		pc := &oapi.NetPeering{}
		for _, c := range result.NetPeerings {
			if rs.Primary.ID == c.NetPeeringId {
				pc = &c
			}
		}

		if pc == nil {
			// not found
			return nil
		}

		if pc.State.Name != "" {
			if pc.State.Name == "deleted" {
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

func testAccCheckOutscaleOAPILinPeeringConnectionExists(n string, connection *oapi.NetPeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Net Peering ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		var resp *oapi.POST_ReadNetPeeringsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadNetPeerings(oapi.ReadNetPeeringsRequest{
				Filters: oapi.FiltersNetPeering{NetPeeringIds: []string{rs.Primary.ID}},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("Status Code: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("Status Code: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("Status: 500, %s", utils.ToJSONString(resp.Code500))
			}
			return fmt.Errorf("Error reading Net Peering details: %s", errString)
		}

		result := resp.OK

		if len(result.NetPeerings) == 0 {
			return fmt.Errorf("Net Peering could not be found")
		}

		*connection = result.NetPeerings[0]

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
// func testAccCheckOutscaleOAPILinPeeringConnectionOptions(n, block string, options *oapi.NetPeeringOptionsDescription) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", n)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No Net Peering ID is set")
// 		}

// 		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

// 		var resp *oapi.ReadNetPeeringsOutput
// 		var err error
// 		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			resp, err = conn.VM.ReadNetPeerings(
// 				&oapi.ReadNetPeeringsRequest{
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
