package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIInternetServiceLink_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
	var conf oapi.InternetService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIInternetServiceLinkDettached,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIInternetServiceLinkConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIInternetServiceLinkExists("outscale_internet_service_link.link", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIInternetServiceLinkExists(n string, res *oapi.InternetService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet service link id is set")
		}

		var resp *oapi.POST_ReadInternetServicesResponses
		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		id := rs.Primary.Attributes["internet_service_id"]

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadInternetServices(oapi.ReadInternetServicesRequest{
				Filters: oapi.FiltersInternetService{InternetServiceIds: []string{id}},
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

			return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", errString)
		}

		result := resp.OK

		if len(result.InternetServices) != 1 ||
			result.InternetServices[0].InternetServiceId != id {
			return fmt.Errorf("Internet Service not found")
		}

		*res = result.InternetServices[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIInternetServiceLinkDettached(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_internet_service_link" {
			continue
		}

		id := rs.Primary.Attributes["internet_gateway_id"]

		// Try to find an internet service
		var resp *oapi.POST_ReadInternetServicesResponses
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadInternetServices(oapi.ReadInternetServicesRequest{
				Filters: oapi.FiltersInternetService{InternetServiceIds: []string{id}},
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

			return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", errString)
		}

		result := resp.OK

		if resp == nil {
			return nil
		}

		if err == nil {
			if len(result.InternetServices) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidInternetGateway.NotFound" {
			return err
		}
	}

	return nil
}

const testAccOutscaleOAPIInternetServiceLinkConfig = `
resource "outscale_internet_service" "gateway" {}

resource "outscale_net" "vpc" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_internet_service_link" "link" {
	net_id = "${outscale_net.vpc.id}"
	internet_service_id = "${outscale_internet_service.gateway.id}"
}
`
