package outscale

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	var conf oscgo.InternetService

	resourceName := "outscale_internet_service.gateway"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOSCAPIInternetServiceDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIInternetServiceConfig("Terraform_IGW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOSCAPIInternetServiceExists(resourceName, &conf),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "Terraform_IGW"),
				),
			},
			{
				Config: testAccOutscaleOAPIInternetServiceConfig("Terraform_IGW2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOSCAPIInternetServiceExists("outscale_internet_service.gateway", &conf),
					resource.TestCheckResourceAttrSet(resourceName, "tags.#"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.key", "Name"),
					resource.TestCheckResourceAttr(resourceName, "tags.0.value", "Terraform_IGW2"),
				),
			},
			{
				Config: testAccOutscaleOAPIInternetServiceWithoutTags(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOSCAPIInternetServiceExists("outscale_internet_service.gateway", &conf),
					resource.TestCheckNoResourceAttr(resourceName, "tags.#"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOSCAPIInternetServiceExists(n string, res *oscgo.InternetService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No internet gateway id is set")
		}
		var resp oscgo.ReadInternetServicesResponse
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			r, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), &oscgo.ReadInternetServicesOpts{ReadInternetServicesRequest: optional.NewInterface(oscgo.ReadInternetServicesRequest{
				Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
			})})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			resp = r
			return nil
		})

		var errString string

		if err != nil {
			errString = err.Error()

			return fmt.Errorf("[DEBUG] Error creating Internet Service: %s", errString)
		}

		if err != nil {
			return err
		}
		if len(resp.GetInternetServices()) != 1 || resp.GetInternetServices()[0].GetInternetServiceId() != rs.Primary.ID {
			return fmt.Errorf("Internet Gateway not found")
		}

		*res = resp.GetInternetServices()[0]

		return nil
	}
}

func testAccCheckOutscaleOSCAPIInternetServiceDestroyed(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_internet_service" {
			continue
		}

		// Try to find an internet gateway
		var resp oscgo.ReadInternetServicesResponse
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error

			r, _, err := conn.InternetServiceApi.ReadInternetServices(context.Background(), &oscgo.ReadInternetServicesOpts{ReadInternetServicesRequest: optional.NewInterface(oscgo.ReadInternetServicesRequest{
				Filters: &oscgo.FiltersInternetService{InternetServiceIds: &[]string{rs.Primary.ID}},
			})})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			resp = r

			return resource.RetryableError(err)
		})

		if err == nil {
			if len(resp.GetInternetServices()) > 0 {
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
