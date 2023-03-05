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

func TestAccOthers_Ca_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_ca.ca_test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleCaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPICaConfig(utils.TestCaPem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resourceName),
				),
			},
			{
				Config: testAccOutscaleOAPICaConfigUpdateDescription(utils.TestCaPem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckOutscaleCaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set")
		}

		var resp oscgo.ReadCasResponse
		var err error
		exists := false
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(oscgo.ReadCasRequest{}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetCas()) == 0 {
			return fmt.Errorf("Ca not found (%s)", rs.Primary.ID)
		}

		for _, ca := range resp.GetCas() {
			if ca.GetCaId() == rs.Primary.ID {
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("Ca not found (%s)", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckOutscaleCaDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_ca" {
			continue
		}
		req := oscgo.ReadCasRequest{}
		req.Filters = &oscgo.FiltersCa{
			CaIds: &[]string{rs.Primary.ID},
		}

		var resp oscgo.ReadCasResponse
		var err error
		exists := false
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("Ca reading (%s)", rs.Primary.ID)
		}

		for _, ca := range resp.GetCas() {
			if ca.GetCaId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("Ca still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleOAPICaConfig(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" { 
   ca_pem        =  %[1]q
   description        = "Ca testacc create"
}
	`, ca_pem)
}

func testAccOutscaleOAPICaConfigUpdateDescription(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" { 
   ca_pem        =  %[1]q
   description        = "Ca testacc update"
}
	`, ca_pem)
}
