package oapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_DataOutscaleCas_basic(t *testing.T) {
	resName := "outscale_ca.ca_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccDataCheckOutscaleCasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleCasConfig(ca_path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resName),
				),
			},
		},
	})
}

func testAccDataCheckOutscaleCasDestroy(s *terraform.State) error {
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
		resp, err := client.ReadCas(context.Background(), req, options.WithRetryTimeout(120*time.Second))
		if err != nil {
			return fmt.Errorf("ca reading (%s)", rs.Primary.ID)
		}

		for _, ca := range ptr.From(resp.Cas) {
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

func testAccDataOutscaleCasConfig(path string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" {
   ca_pem        =  file(%[1]q)
   description   = "Ca testacc create"
}

resource "outscale_ca" "ca_test2" {
   ca_pem        = file(%[1]q)
   description   = "Ca testacc create2"
}

data "outscale_cas" "cas_data" {
   filter {
      name   = "descriptions"
      values = ["Ca testacc create*"]
   }
}`, path)
}
