package outscale

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleCustomerGateway_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var gateway fcu.CustomerGateway
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_client_endpoint.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGatewayConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGateway("outscale_client_endpoint.foo", &gateway),
				),
			},
			{
				Config: testAccCustomerGatewayConfigUpdateTags(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGateway("outscale_client_endpoint.foo", &gateway),
				),
			},
			{
				Config: testAccCustomerGatewayConfigForceReplace(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGateway("outscale_client_endpoint.foo", &gateway),
				),
			},
		},
	})
}

func TestAccOutscaleCustomerGateway_similarAlreadyExists(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var gateway fcu.CustomerGateway
	rInt := acctest.RandInt()
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_client_endpoint.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGatewayConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGateway("outscale_client_endpoint.foo", &gateway),
				),
			},
			{
				Config:      testAccCustomerGatewayConfigIdentical(rInt, rBgpAsn),
				ExpectError: regexp.MustCompile("An existing customer gateway"),
			},
		},
	})
}

func TestAccOutscaleCustomerGateway_disappears(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	rInt := acctest.RandInt()
	rBgpAsn := acctest.RandIntRange(64512, 65534)
	var gateway fcu.CustomerGateway
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomerGatewayConfig(rInt, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCustomerGateway("outscale_client_endpoint.foo", &gateway),
					testAccAWSCustomerGatewayDisappears(&gateway),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAWSCustomerGatewayDisappears(gateway *fcu.CustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		opts := &fcu.DeleteCustomerGatewayInput{
			CustomerGatewayId: gateway.CustomerGatewayId,
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := conn.VM.DeleteCustomerGateway(opts)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		return resource.Retry(40*time.Minute, func() *resource.RetryError {
			opts := &fcu.DescribeCustomerGatewaysInput{
				CustomerGatewayIds: []*string{gateway.CustomerGatewayId},
			}
			resp, err := conn.VM.DescribeCustomerGateways(opts)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("Error retrieving Customer Gateway: %s", err))
			}
			if *resp.CustomerGateways[0].State == "deleted" {
				return nil
			}
			return resource.RetryableError(fmt.Errorf(
				"Waiting for Customer Gateway: %v", gateway.CustomerGatewayId))
		})
	}
}

func testAccCheckCustomerGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_customer_endpoint" {
			continue
		}

		gatewayFilter := &fcu.Filter{
			Name:   aws.String("customer-gateway-id"),
			Values: []*string{aws.String(rs.Primary.ID)},
		}

		resp, err := conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
			Filters: []*fcu.Filter{gatewayFilter},
		})

		if strings.Contains(fmt.Sprint(err), "InvalidCustomerGatewayID.NotFound") {
			continue
		}

		if err == nil {
			if len(resp.CustomerGateways) > 0 {
				return fmt.Errorf("Customer gateway still exists: %v", resp.CustomerGateways)
			}

			if *resp.CustomerGateways[0].State == "deleted" {
				continue
			}
		}

		return err
	}

	return nil
}

func testAccCheckCustomerGateway(gatewayResource string, cgw *fcu.CustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[gatewayResource]
		if !ok {
			return fmt.Errorf("Not found: %s", gatewayResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		gateway, ok := s.RootModule().Resources[gatewayResource]
		if !ok {
			return fmt.Errorf("Not found: %s", gatewayResource)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		gatewayFilter := &fcu.Filter{
			Name:   aws.String("customer-gateway-id"),
			Values: []*string{aws.String(gateway.Primary.ID)},
		}

		var resp *fcu.DescribeCustomerGatewaysOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeCustomerGateways(&fcu.DescribeCustomerGatewaysInput{
				Filters: []*fcu.Filter{gatewayFilter},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		respGateway := resp.CustomerGateways[0]
		*cgw = *respGateway

		return nil
	}
}

func testAccCustomerGatewayConfig(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_client_endpoint" "foo" {
			bgp_asn = %d
			ip_address = "172.0.0.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-%d"
			}
		}
		`, rBgpAsn, rInt)
}

func testAccCustomerGatewayConfigIdentical(randInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_client_endpoint" "foo" {
			bgp_asn = %d
			ip_address = "172.0.0.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-%d"
			}
		}
		resource "outscale_client_endpoint" "identical" {
			bgp_asn = %d
			ip_address = "172.0.0.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-identical-%d"
			}
		}
		`, rBgpAsn, randInt, rBgpAsn, randInt)
}

// Add the Another: "tag" tag.
func testAccCustomerGatewayConfigUpdateTags(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
	resource "outscale_client_endpoint" "foo" {
		bgp_asn = %d
		ip_address = "172.0.0.1"
		type = "ipsec.1"
		tag {
			Name = "foo-gateway-%d"
			Another = "tag"
		}
	}
	`, rBgpAsn, rInt)
}

// Change the ip_address.
func testAccCustomerGatewayConfigForceReplace(rInt, rBgpAsn int) string {
	return fmt.Sprintf(`
		resource "outscale_client_endpoint" "foo" {
			bgp_asn = %d
			ip_address = "172.10.10.1"
			type = "ipsec.1"
			tag {
				Name = "foo-gateway-%d"
				Another = "tag"
			}
		}
		`, rBgpAsn, rInt)
}
