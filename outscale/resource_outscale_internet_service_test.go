package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccInternetService_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInternetServiceDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccInternetServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInternetServiceExists("outscale_internet_service.outscale_internet_service"),
				),
			},
			{
				Config: testAccInternetServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInternetServiceExists("outscale_internet_service.outscale_internet_service"),
				),
			},
		},
	})
}

func testAccCheckInternetServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*Client).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}

		filterReq := oscgo.ReadInternetServicesRequest{
			Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadInternetServicesResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetInternetServices()) < 1 {
			return fmt.Errorf("Internet Service Link not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckInternetServiceDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_internet_service_link" {
			continue
		}

		filterReq := oscgo.ReadInternetServicesRequest{
			Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadInternetServicesResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(filterReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp.StatusCode, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetInternetServices()) > 0 {
			return fmt.Errorf("Internet Service Link still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccInternetServiceConfig() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-internet-service-rs"
			}
		}

		resource "outscale_internet_service" "outscale_internet_service" {
			tags {
				key = "Name"
				value = "testacc-internet-service"
			}
		}
	`
}
