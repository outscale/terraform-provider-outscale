package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_Ca_Basic(t *testing.T) {
	resourceName := "outscale_ca.ca_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
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
		Steps: testacc.FrameworkMigrationTestSteps("1.3.1", testAccOutscaleCaConfig(testAccCertPath)),
	})
}

func testAccCheckOutscaleCaExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := testacc.ConfiguredClient.OSC

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		exists := false
		resp, err := client.ReadCas(context.Background(), osc.ReadCasRequest{}, options.WithRetryTimeout(DefaultTimeout))

		if err != nil || resp.Cas == nil || len(*resp.Cas) == 0 {
			return fmt.Errorf("ca not found (%s)", rs.Primary.ID)
		}

		for _, ca := range *resp.Cas {
			if *ca.CaId == rs.Primary.ID {
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
	client := testacc.ConfiguredClient.OSC

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_ca" {
			continue
		}
		req := osc.ReadCasRequest{}
		req.Filters = &osc.FiltersCa{
			CaIds: &[]string{rs.Primary.ID},
		}

		exists := false
		resp, err := client.ReadCas(context.Background(), req, options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			return fmt.Errorf("ca reading (%s)", rs.Primary.ID)
		}

		for _, ca := range *resp.Cas {
			if *ca.CaId == rs.Primary.ID {
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
