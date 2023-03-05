package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_DataOutscaleCas_basic(t *testing.T) {
	resName := "outscale_ca.ca_test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDataCheckOutscaleCasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOutscaleOAPICasConfig(utils.TestCaPem),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleCaExists(resName),
				),
			},
		},
	})
}

func testAccDataCheckOutscaleCasDestroy(s *terraform.State) error {
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

func testAccDataOutscaleOAPICasConfig(ca_pem string) string {
	return fmt.Sprintf(`
resource "outscale_ca" "ca_test" {
   ca_pem        =  %[1]q
   description        = "Ca testacc create"
}

resource "outscale_ca" "ca_test2" { 
   ca_pem        =  %[1]q
   description        = "Ca testacc create2"
}

data "outscale_cas" "cas_data" {
   filter {
      name   = "ca_ids"
      values = [outscale_ca.ca_test.id]
   }
   filter {
      name   = "description"
      values = ["Ca testacc create2"]
   }
}
`, ca_pem)
}
