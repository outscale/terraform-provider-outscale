package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPISubNet_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
	var conf oapi.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckOutscaleLinDestroyed, //TODO: fix net destroy test
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPISubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISubNetExists("outscale_subnet.subnet", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISubNetExists(n string, res *oapi.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet id is set")
		}
		var resp *oapi.POST_ReadSubnetsResponses
		conn := testAccProvider.Meta().(*OutscaleClient)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.OAPI.POST_ReadSubnets(oapi.ReadSubnetsRequest{
				Filters: oapi.FiltersSubnet{SubnetIds: []string{rs.Primary.ID}},
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
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("[DEBUG] Error reading Subnet (%s)", errString)
		}

		response := resp.OK

		if len(response.Subnets) != 1 ||
			response.Subnets[0].SubnetId != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*res = response.Subnets[0]

		return nil
	}
}

func testAccCheckOutscaleOAPISubNetDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_subnet" {
			continue
		}

		// Try to find an internet gateway
		var resp *oapi.POST_ReadSubnetsResponses
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.OAPI.POST_ReadSubnets(oapi.ReadSubnetsRequest{
				Filters: oapi.FiltersSubnet{SubnetIds: []string{rs.Primary.ID}},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if resp.OK == nil {
			return nil
		}

		if err == nil {
			if len(resp.OK.Subnets) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		// Verify the error is what we want
		if !strings.Contains(fmt.Sprintf("%s", err), "InvalidSubnet.NotFound") {
			return err
		}
	}

	return nil
}

const testAccOutscaleOAPISubnetConfig = `
resource "outscale_net" "net" {
	ip_range = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	ip_range = "10.0.0.0/16"
	subregion_name = "eu-west-2a"
	net_id = "${outscale_net.net.id}"

	tags = {
		key = "name"
		value = "terraform-subnet"
	 }
}

`
