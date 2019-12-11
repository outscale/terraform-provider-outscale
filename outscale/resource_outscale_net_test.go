package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPILin_basic(t *testing.T) {
	var conf1 oscgo.Net
	var conf2 oscgo.Net

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckOutscaleLinDestroyed, // we need to create the destroyed test case
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPILinExists("outscale_net.vpc.0", &conf1),
					testAccCheckOutscaleOAPILinExists("outscale_net.vpc.1", &conf2),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc.0", "ip_range", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(
						"outscale_net.vpc.1", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPILin_UpdateTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPILinConfigUpdateTags("Terraform_net"),
				Check:  resource.ComposeTestCheckFunc(),
			},
			resource.TestStep{
				Config: testAccOutscaleOAPILinConfigUpdateTags("Terraform_net2"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckOutscaleOAPILinExists(n string, res *oscgo.Net) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp oscgo.ReadNetsResponse
		conn := testAccProvider.Meta().(*OutscaleClient)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, _, err = conn.OSCAPI.NetApi.ReadNets(context.Background(), &oscgo.ReadNetsOpts{ReadNetsRequest: optional.NewInterface(oscgo.ReadNetsRequest{
				Filters: &oscgo.FiltersNet{NetIds: &[]string{rs.Primary.ID}},
			})})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return err
		}

		if len(resp.GetNets()) != 1 ||
			resp.GetNets()[0].GetNetId() != rs.Primary.ID {
			return fmt.Errorf("Net not found")
		}

		*res = resp.GetNets()[0]

		return nil
	}
}

//Missing on Swagger Spec
// func testAccCheckOutscaleOAPILinDestroyed(s *terraform.State) error {
// 	conn := testAccProvider.Meta().(*OutscaleClient)

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "outscale_net" {
// 			continue
// 		}

// 		// Try to find an internet gateway
// 		var resp *oscgo.ReadGate
// 		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			var err error
// 			resp, err = conn.FCU.VM.DescribeInternetGateways(&fcu.DescribeInternetGatewaysInput{
// 				InternetGatewayIds: []*string{aws.String(rs.Primary.ID)},
// 			})

// 			if err != nil {
// 				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
// 					return resource.RetryableError(err)
// 				}
// 				return resource.NonRetryableError(err)
// 			}

// 			return resource.RetryableError(err)
// 		})

// 		if resp == nil {
// 			return nil
// 		}

// 		if err == nil {
// 			if len(resp.InternetGateways) > 0 {
// 				return fmt.Errorf("still exist")
// 			}
// 			return nil
// 		}

// 		// Verify the error is what we want
// 		ec2err, ok := err.(awserr.Error)
// 		if !ok {
// 			return err
// 		}
// 		if ec2err.Code() != "InvalidVPC.NotFound" {
// 			return err
// 		}
// 	}

// 	return nil
// }

const testAccOutscaleOAPILinConfig = `
	resource "outscale_net" "vpc" {
		ip_range = "10.0.0.0/16"
		count = 2

		tags {
			key = "Name" 
			value = "outscale_net"
		}	
	}
`

func testAccOutscaleOAPILinConfigUpdateTags(value string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" { 
		ip_range = "10.0.0.0/16"
		tags { 
			key = "name" 
			value = "%s"
		}
	   }
`, value)
}
