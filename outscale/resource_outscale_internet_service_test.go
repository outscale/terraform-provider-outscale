package outscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIInternetService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleInternetServiceDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIInternetServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleInternetServiceExists("outscale_internet_service.outscale_internet_service"),
				),
			},
			{
				Config: testAccOutscaleOAPIInternetServiceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleInternetServiceExists("outscale_internet_service.outscale_internet_service"),
				),
			},
		},
	})
}

func testAccCheckOutscaleInternetServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}

		filterReq := &oscgo.ReadInternetServicesOpts{
			ReadInternetServicesRequest: optional.NewInterface(oscgo.ReadInternetServicesRequest{
				Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
			}),
		}

		resp, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), filterReq)
		if err != nil || len(resp.GetInternetServices()) < 1 {
			return fmt.Errorf("Internet Service Link not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckOutscaleInternetServiceDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_internet_service_link" {
			continue
		}

		filterReq := &oscgo.ReadInternetServicesOpts{
			ReadInternetServicesRequest: optional.NewInterface(oscgo.ReadInternetServicesRequest{
				Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
			}),
		}

		resp, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), filterReq)
		if err != nil || len(resp.GetInternetServices()) > 0 {
			return fmt.Errorf("Internet Service Link still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleOAPIInternetServiceConfig() string {
	return `
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_internet_service" "outscale_internet_service" {}
	`
}

func testAccOutscaleOAPIInternetServiceWithoutTags() string {
	return `resource "outscale_internet_service" "gateway" {}`
}
