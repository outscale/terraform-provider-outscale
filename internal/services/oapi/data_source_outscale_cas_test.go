package oapi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_DataOutscaleCas_basic(t *testing.T) {
	resName := "outscale_ca.ca_test"
	ca_path := testAccCertPath

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
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
	conn := testacc.ConfiguredClient.OSCAPI

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
