package outscale

// import (
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/awserr"
// 	"github.com/hashicorp/terraform/helper/resource"
// 	"github.com/hashicorp/terraform/terraform"
// 	"github.com/terraform-providers/terraform-provider-outscale/osc/lbu"
// )

// func TestAccOutscaleOAPILBU_basic(t *testing.T) {
// 	o := os.Getenv("OUTSCALE_OAPI")

// 	oapi, err := strconv.ParseBool(o)
// 	if err != nil {
// 		oapi = false
// 	}

// 	if !oapi {
// 		t.Skip()
// 	}

// 	var conf lbu.LoadBalancerDescription

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:      func() { testAccPreCheck(t) },
// 		IDRefreshName: "outscale_load_balancer.bar",
// 		Providers:     testAccProviders,
// 		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccDSOutscaleOAPILBUConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckOutscaleOAPILBUExists("outscale_load_balancer.bar", &conf),
// 					testAccCheckOutscaleOAPILBUAttributes(&conf),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.#", "2"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.0", "eu-west-2a"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "sub_region_name.1", "eu-west-2b"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.backend_port", "8000"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.backend_protocol", "http"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.load_balancer_port", "80"),
// 					resource.TestCheckResourceAttr(
// 						"data.outscale_load_balancer.test", "listener.0.load_balancer_protocol", "http"),
// 				)},
// 		},
// 	})
// }

// func testAccCheckOutscaleOAPILBUDestroy(s *terraform.State) error {
// 	conn := testAccProvider.Meta().(*OutscaleClient).LBU

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "outscale_load_balancer" {
// 			continue
// 		}

// 		var err error
// 		var describe *lbu.DescribeLoadBalancersOutput
// 		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
// 				LoadBalancerNames: &lbu.LoadBalancerNamesMember{Member: []*string{aws.String(rs.Primary.ID)}},
// 			})

// 			if err != nil {
// 				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
// 					return resource.RetryableError(
// 						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
// 				}
// 				return resource.NonRetryableError(err)
// 			}
// 			return nil
// 		})

// 		if err == nil {
// 			if len(describe.LoadBalancerDescriptions) != 0 &&
// 				*describe.LoadBalancerDescriptions[0].LoadBalancerName == rs.Primary.ID {
// 				return fmt.Errorf("LBU still exists")
// 			}
// 		}

// 		// Verify the error
// 		providerErr, ok := err.(awserr.Error)
// 		if !ok {
// 			return err
// 		}

// 		if providerErr.Code() != "LoadBalancerNotFound" {
// 			return fmt.Errorf("Unexpected error: %s", err)
// 		}
// 	}

// 	return nil
// }

// func testAccCheckOutscaleOAPILBUAttributes(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
// 		azs := make([]string, 0, len(conf.AvailabilityZones.Member))
// 		for _, x := range conf.AvailabilityZones.Member {
// 			azs = append(azs, *x)
// 		}
// 		sort.StringSlice(azs).Sort()
// 		if !reflect.DeepEqual(azs, zones) {
// 			return fmt.Errorf("bad availability_zones_member")
// 		}

// 		l := lbu.Listener{
// 			InstancePort:     aws.Int64(int64(8000)),
// 			InstanceProtocol: aws.String("HTTP"),
// 			LoadBalancerPort: aws.Int64(int64(80)),
// 			Protocol:         aws.String("HTTP"),
// 		}

// 		if !reflect.DeepEqual(conf.ListenerDescriptions[0].Listener, &l) {
// 			return fmt.Errorf(
// 				"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
// 				conf.ListenerDescriptions[0].Listener,
// 				l)
// 		}

// 		if *conf.DNSName == "" {
// 			return fmt.Errorf("empty dns_name")
// 		}

// 		return nil
// 	}
// }

// func testAccCheckOutscaleOAPILBUAttributesHealthCheck(conf *lbu.LoadBalancerDescription) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		zones := []string{"eu-west-2a", "eu-west-2b", "eu-west-2c"}
// 		azs := make([]string, 0, len(conf.AvailabilityZones.Member))
// 		for _, x := range conf.AvailabilityZones.Member {
// 			azs = append(azs, *x)
// 		}
// 		sort.StringSlice(azs).Sort()
// 		if !reflect.DeepEqual(azs, zones) {
// 			return fmt.Errorf("bad availability_zones_member")
// 		}

// 		check := &lbu.HealthCheck{
// 			Timeout:            aws.Int64(int64(30)),
// 			UnhealthyThreshold: aws.Int64(int64(5)),
// 			HealthyThreshold:   aws.Int64(int64(5)),
// 			Interval:           aws.Int64(int64(60)),
// 			Target:             aws.String("HTTP:8000/"),
// 		}

// 		if !reflect.DeepEqual(conf.HealthCheck, check) {
// 			return fmt.Errorf(
// 				"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
// 				conf.HealthCheck,
// 				check)
// 		}

// 		if *conf.DNSName == "" {
// 			return fmt.Errorf("empty dns_name")
// 		}

// 		return nil
// 	}
// }

// func testAccCheckOutscaleOAPILBUExists(n string, res *lbu.LoadBalancerDescription) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", n)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No LBU ID is set")
// 		}

// 		conn := testAccProvider.Meta().(*OutscaleClient).LBU

// 		var err error
// 		var describe *lbu.DescribeLoadBalancersOutput
// 		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
// 			describe, err = conn.API.DescribeLoadBalancers(&lbu.DescribeLoadBalancersInput{
// 				LoadBalancerNames: &lbu.LoadBalancerNamesMember{Member: []*string{aws.String(rs.Primary.ID)}},
// 			})

// 			if err != nil {
// 				if strings.Contains(fmt.Sprint(err), "CertificateNotFound") {
// 					return resource.RetryableError(
// 						fmt.Errorf("[WARN] Error creating ELB Listener with SSL Cert, retrying: %s", err))
// 				}
// 				return resource.NonRetryableError(err)
// 			}
// 			return nil
// 		})

// 		if err != nil {
// 			return err
// 		}

// 		if len(describe.LoadBalancerDescriptions) != 1 ||
// 			*describe.LoadBalancerDescriptions[0].LoadBalancerName != rs.Primary.ID {
// 			return fmt.Errorf("LBU not found")
// 		}

// 		*res = *describe.LoadBalancerDescriptions[0]

// 		// Confirm source_security_group_id for ELBs in a VPC
// 		// 	See https://github.com/hashicorp/terraform/pull/3780
// 		if res.VPCId != nil {
// 			sgid := rs.Primary.Attributes["source_security_group_id"]
// 			if sgid == "" {
// 				return fmt.Errorf("Expected to find source_security_group_id for LBU, but was empty")
// 			}
// 		}

// 		return nil
// 	}
// }

// const testAccOutscaleOAPILBUConfig = `
// resource "outscale_load_balancer" "bar" {
//   sub_region_name = ["eu-west-2a", "eu-west-2b"]
// 	load_balancer_name               = "foobar-terraform-elb"
//   listener {
//     backend_port = 8000
//     backend_protocol = "http"
//     load_balancer_port = 80
//     // Protocol should be case insensitive
//     load_balancer_protocol = "http"
//   }

// 	tag {
// 		bar = "baz"
// 	}
// }
// `
