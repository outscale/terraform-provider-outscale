package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/outscale/osc-go/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIPublicIP_basic(t *testing.T) {
	var conf oapi.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip" {
			continue
		}
		//Missing on Swagger Spec
		// if strings.Contains(rs.Primary.ID, "reservation") {
		// 	req := oapi.ReadPublicIpsRequest{
		// 		Filters: oapi.FiltersPublicIpcIp{
		// 			ReservationIds: []string{rs.Primary.ID},
		// 		},
		// 	}

		// 	var describe *oapi.ReadPublicIpsResponse
		// 	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		// 		var err error
		// 		resp, err := conn.OAPI.POST_ReadPublicIps(req)
		// 		describe = resp.OK
		// 		return resource.RetryableError(err)
		// 	})

		// 	if err != nil {
		// 		// Verify the error is what we want
		// 		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
		// 			return nil
		// 		}

		// 		return err
		// 	}

		// 	if len(describe.PublicIps) > 0 {
		// 		return fmt.Errorf("still exists")
		// 	}
		// } else {
		req := oapi.ReadPublicIpsRequest{
			Filters: oapi.FiltersPublicIp{
				PublicIps: []string{rs.Primary.ID},
			},
		}

		var describe *oapi.ReadPublicIpsResponse
		err := resource.Retry(60*time.Second, func() *resource.RetryError {
			var err error
			resp, err := conn.OAPI.POST_ReadPublicIps(req)
			describe = resp.OK
			return resource.RetryableError(err)
		})

		if err != nil {
			// Verify the error is what we want
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
				return nil
			}

			return err
		}

		if len(describe.PublicIps) > 0 {
			return fmt.Errorf("still exists")
		}
		//}
	}

	return nil
}

func testAccCheckOutscaleOAPIPublicIPAttributes(conf *oapi.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conf.PublicIp == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIPublicIPExists(n string, res *oapi.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		req := oapi.ReadPublicIpsRequest{
			Filters: oapi.FiltersPublicIp{
				PublicIpIds: []string{rs.Primary.ID},
			},
		}

		var describe *oapi.ReadPublicIpsResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			resp, err := conn.OAPI.POST_ReadPublicIps(req)

			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			describe = resp.OK

			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
				return nil
			}

			return err
		}

		if err != nil {

			// Verify the error is what we want
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
				return nil
			}

			return err
		}

		if len(describe.PublicIps) != 1 ||
			describe.PublicIps[0].PublicIpId != rs.Primary.ID {
			return fmt.Errorf("PublicIP not found")
		}
		*res = describe.PublicIps[0]
		//}

		return nil
	}
}

const testAccOutscaleOAPIPublicIPConfig = `
resource "outscale_public_ip" "bar" {
	tags {
		key = "Name"
		value = "public_ip_test"
	}
}
`
