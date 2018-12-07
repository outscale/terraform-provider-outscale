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

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIInternetService_basic(t *testing.T) {
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
		CheckDestroy: testAccCheckOutscaleOAPIInternetServiceDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIInternetServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIInternetServiceExists("outscale_internet_service.gateway", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIInternetServiceExists(n string, res *oapi.InternetService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp *oapi.POST_ReadInternetServicesResponses
		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadInternetServices(oapi.ReadInternetServicesRequest{
				Filters: oapi.FiltersInternetService{InternetServiceIds: []string{rs.Primary.ID}},
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

			return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", errString)
		}

		result := resp.OK

		if err != nil {
			return err
		}
		if len(result.InternetServices) != 1 ||
			result.InternetServices[0].InternetServiceId != rs.Primary.ID {
			return fmt.Errorf("Internet Gateway not found")
		}

		*res = result.InternetServices[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIInternetServiceDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_internet_service" {
			continue
		}

		// Try to find an internet gateway
		var resp *oapi.POST_ReadInternetServicesResponses
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error

			resp, err = conn.POST_ReadInternetServices(oapi.ReadInternetServicesRequest{
				Filters: oapi.FiltersInternetService{InternetServiceIds: []string{rs.Primary.ID}},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		})

		if resp.OK == nil {
			return nil
		}

		result := resp.OK

		if err == nil {
			if len(result.InternetServices) > 0 {
				return fmt.Errorf("still exist")
			}
			return nil
		}

		if !strings.Contains(fmt.Sprintf("%s", err), "InvalidInternetService.NotFound") {
			return err
		}
	}

	return nil
}

const testAccOutscaleOAPIInternetServiceConfig = `
resource "outscale_internet_service" "gateway" {}
`
