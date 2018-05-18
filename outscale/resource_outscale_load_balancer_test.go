package outscale

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
)

func TestAccOutscaleLBU_basic(t *testing.T) {
	var conf lbu.LoadBalancerDescription

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_load_balancer.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleLBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLBUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleLBUExists("outscale_load_balancer.bar", &conf),
					testAccCheckOutscaleLBUAttributes(&conf),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.#", "2"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.0", "eu-west-2a"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "availability_zones_member.1", "eu-west-2b"),
					// resource.TestCheckResourceAttr(
					// 	"outscale_load_balancer.bar", "subnets.#", "3"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.instance_port", "8000"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.instance_protocol", "http"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.load_balancer_port", "80"),
					resource.TestCheckResourceAttr(
						"outscale_load_balancer.bar", "listeners_member.0.protocol", "http"),
				)},
		},
	})
}

func testAccCheckOutscaleLBUDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).LBU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_load_balancer" {
			continue
		}

		var err error
		var describe *lbu.DescribeLoadBalancersOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err == nil {
			if len(describe.LoadBalancerDescriptions) != 0 &&
				*describe.LoadBalancerDescriptions[0].LoadBalancerName == rs.Primary.ID {
				return fmt.Errorf("LBU still exists")
			}
		}

		// Verify the error
		providerErr, ok := err.(awserr.Error)
		if !ok {
			return err
		}

		if providerErr.Code() != "LoadBalancerNotFound" {
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}

	return nil
}

func testAccCheckOutscaleLBUAttributes(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
		azs := make([]string, 0, len(conf.AvailabilityZones))
		for _, x := range conf.AvailabilityZones {
			azs = append(azs, *x)
		}
		sort.StringSlice(azs).Sort()
		if !reflect.DeepEqual(azs, zones) {
			return fmt.Errorf("bad availability_zones_member")
		}

		if *conf.DNSName == "" {
			return fmt.Errorf("empty dns_name")
		}

		return nil
	}
}

func testAccCheckOutscaleLBUAttributesHealthCheck(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
		azs := make([]string, 0, len(conf.AvailabilityZones))
		for _, x := range conf.AvailabilityZones {
			azs = append(azs, *x)
		}
		sort.StringSlice(azs).Sort()
		if !reflect.DeepEqual(azs, zones) {
			return fmt.Errorf("bad availability_zones_member")
		}

		if *conf.DNSName == "" {
			return fmt.Errorf("empty dns_name")
		}

		return nil
	}
}

func testAccCheckOutscaleLBUExists(n string, res *lbu.LoadBalancerDescription) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LBU ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).LBU

		var err error
		var describe *lbu.DescribeLoadBalancersOutput
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
				LoadBalancerNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
					return resource.RetryableError(
						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(describe.LoadBalancerDescriptions) != 1 ||
			*describe.LoadBalancerDescriptions[0].LoadBalancerName != rs.Primary.ID {
			return fmt.Errorf("LBU not found")
		}

		*res = *describe.LoadBalancerDescriptions[0]

		// Confirm source_security_group_id for ELBs in a VPC
		// 	See https://github.com/hashicorp/terraform/pull/3780
		if res.VPCId != nil {
			sgid := rs.Primary.Attributes["source_security_group_id"]
			if sgid == "" {
				return fmt.Errorf("Expected to find source_security_group_id for LBU, but was empty")
			}
		}

		return nil
	}
}

const testAccOutscaleLBUConfig = `
resource "outscale_load_balancer" "bar" {
  availability_zones_member = ["eu-west-2a", "eu-west-2b"]
	load_balancer_name               = "foobar-terraform-elb"
  listeners_member {
    instance_port = 8000
    instance_protocol = "http"
    load_balancer_port = 80
    // Protocol should be case insensitive
    protocol = "http"
  }

	tag {
		bar = "baz"
	}

}
`
