package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccAWSDHCPOptions_basic(t *testing.T) {
	var d fcu.DhcpOptions
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDHCPOptionsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDHCPOptionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDHCPOptionsExists("outscale_dhcp_options.foo", &d),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "dhcp_configuration_set", "service.consul"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "dhcp_configuration_set.key", "127.0.0.1"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "dhcp_configuration_set.value_set.value", "10.0.0.2"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "dhcp_options_id.0", "127.0.0.1"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "tag_set.key", "127.0.0.1"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "tag_set.value", "2"),
					resource.TestCheckResourceAttr("outscale_dhcp_options.foo", "request_id", "foo-name"),
				),
			},
		},
	})
}

func TestAccOutscaleDHCPOptions_deleteOptions(t *testing.T) {
	var d fcu.DhcpOptions
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDHCPOptionsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDHCPOptionsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDHCPOptionsExists("outscale_dhcp_options.foo", &d),
					testAccCheckDHCPOptionsDelete("outscale_dhcp_options.foo"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDHCPOptionsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).FCU

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_dhcp_options" {
			continue
		}

		// Try to find the resource
		resp, err := conn.VM.DescribeDhcpOptions(&fcu.DescribeDhcpOptionsInput{
			DhcpOptionsIds: []*string{
				aws.String(rs.Primary.ID),
			},
		})
		if ae, ok := err.(awserr.Error); ok && ae.Code() == "InvalidDhcpOptionID.NotFound" {
			continue
		}

		if strings.Contains(fmt.Sprint(err), "InvalidDhcpOptionID.NotFound") {
		}

		if err == nil {
			if len(resp.DhcpOptions) > 0 {
				return fmt.Errorf("still exists")
			}

			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidDhcpOptionsID.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckDHCPOptionsExists(n string, d *fcu.DhcpOptions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		// //	resp, err := conn.VM.DescribeDhcpOptions(&fcu.DescribeDhcpOptionsInput{
		// 		DhcpOptionsIds: []*string{
		// 			aws.String(rs.Primary.ID),
		// 		},
		// 	})

		var resp *fcu.DescribeDhcpOptionsOutput

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.DescribeDhcpOptions(&fcu.DescribeDhcpOptionsInput)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {

			return fmt.Errorf("DHCP Options not found: %s", err)
		}

		// if err != nil {
		// 	return err
		// }
		// if len(resp.DhcpOptions) == 0 {
		// 	return fmt.Errorf("DHCP Options not found")
		// }

		*d = *resp.DhcpOptions[0]

		return nil
	}
}

func testAccCheckDHCPOptionsDelete(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU
		_, err := conn.VM.DeleteDhcpOptions(&fcu.DeleteDhcpOptionsInput{
			DhcpOptionsId: aws.String(rs.Primary.ID),
		})

		return err
	}
}

const testAccDHCPOptionsConfig = `
resource "outscale_dhcp_options" "foo" {
	dhcp_configuration_set {
	 key = "service.consul"
	 value_set {
		 value = ["127.0.0.1"]
	 }
	} 
	tags {
		Key = "foo-name"
	}
}
`
