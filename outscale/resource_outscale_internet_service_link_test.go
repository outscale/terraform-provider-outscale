package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithInternetServiceLink_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOSCAPIInternetServiceLinkDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOSCAPIInternetServiceLinkExists("outscale_internet_service_link.outscale_internet_service_link"),
				),
			},
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOSCAPIInternetServiceLinkExists("outscale_internet_service_link.outscale_internet_service_link"),
				),
			},
		},
	})
}

func TestAccNet_WithImportInternetServiceLink_Basic(t *testing.T) {
	resourceName := "outscale_internet_service_link.outscale_internet_service_link"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOSCAPIInternetServiceLinkDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleInternetServiceLinkConfig(),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleInternetServiceLinkImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccCheckOutscaleInternetServiceLinkImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckOutscaleOSCAPIInternetServiceLinkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

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
				return utils.CheckThrottling(httpResp, err)
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

func testAccCheckOutscaleOSCAPIInternetServiceLinkDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

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
				return utils.CheckThrottling(httpResp, err)
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

func testAccOutscaleInternetServiceLinkConfig() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-internet-service-link-rs"
			}
		}

		resource "outscale_internet_service" "outscale_internet_service" {}

		resource "outscale_internet_service_link" "outscale_internet_service_link" {
			net_id              = outscale_net.outscale_net.net_id
			internet_service_id = outscale_internet_service.outscale_internet_service.id
		}
	`
}
