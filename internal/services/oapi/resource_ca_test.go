package oapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_Ca_Basic(t *testing.T) {
	resourceName := "outscale_ca.ca_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscaleCaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleCaConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resourceName),
				),
			},
			{
				Config: testAccOutscaleCaConfigUpdateDescription(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resourceName),
				),
			},
			testacc.ImportStep(resourceName, append(testacc.DefaultIgnores(), "ca_pem")...),
		},
	})
}

func TestAccOthers_Ca_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps:    testacc.FrameworkMigrationTestSteps("1.3.1", testAccOutscaleCaConfig(testAccCertPath)),
	})
}

func testAccCheckOutscaleCaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		conn := testacc.ConfiguredClient.OSCAPI

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		var resp osc.ReadCasResponse
		var err error
		exists := false
		err = retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(osc.ReadCasRequest{}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil || len(resp.GetCas()) == 0 {
			return fmt.Errorf("ca not found (%s)", rs.Primary.ID)
		}

		for _, ca := range resp.GetCas() {
			if ca.GetCaId() == rs.Primary.ID {
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("ca not found (%s)", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckOutscaleCaDestroy(s *terraform.State) error {
	conn := testacc.ConfiguredClient.OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_ca" {
			continue
		}
		req := osc.ReadCasRequest{}
		req.Filters = &osc.FiltersCa{
			CaIds: &[]string{rs.Primary.ID},
		}

		var resp osc.ReadCasResponse
		var err error
		exists := false
		err = retry.Retry(120*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.CaApi.ReadCas(context.Background()).ReadCasRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("ca reading (%s)", rs.Primary.ID)
		}

		for _, ca := range resp.GetCas() {
			if ca.GetCaId() == rs.Primary.ID {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("ca still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccOutscaleCaConfig(path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" {
   ca_pem       = file(%q)
   description  = "Ca testacc create"
}`, path)
}

func testAccOutscaleCaConfigUpdateDescription(path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" {
   ca_pem       = file(%q)
   description  = "Ca testacc update"
}`, path)
}
