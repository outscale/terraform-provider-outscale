package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIInternetService_basic(t *testing.T) {
	var conf oapi.InternetService

	resourceName := "outscale_internet_service.gateway"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIInternetServiceDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIInternetServiceConfig("Terraform_IGW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIInternetServiceExists(resourceName, &conf),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "Terraform_IGW"),
				),
			},
			{
				Config: testAccOutscaleOAPIInternetServiceConfig("Terraform_IGW2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIInternetServiceExists("outscale_internet_service.gateway", &conf),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "Terraform_IGW2"),
				),
			},
			{
				Config: testAccOutscaleOAPIInternetServiceWithoutTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIInternetServiceExists("outscale_internet_service.gateway", &conf),
					resource.TestCheckNoResourceAttr(resourceName, "tags.#"),
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

func testAccOutscaleOAPIInternetServiceConfig(value string) string {
	return fmt.Sprintf(`
	resource "outscale_internet_service" "gateway" {
		tags {       
			key   = "Name"     
			value = "%s"       
		}
	}`, value)
}

func testAccOutscaleOAPIInternetServiceWithoutTags() string {
	return `resource "outscale_internet_service" "gateway" {}`
}
