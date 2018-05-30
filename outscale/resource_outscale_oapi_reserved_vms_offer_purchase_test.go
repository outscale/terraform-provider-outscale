package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIReservedVmsOfferPurchase_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	t.Skip()

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIReservedVmsOfferPurchaseEgressConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIReservedVmsOfferPurchaseExists("outscale_reserved_vms_offer_purchase.test"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIReservedVmsOfferPurchaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Group is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		req := &fcu.DescribeReservedInstancesOfferingsInput{
			ReservedInstancesOfferingIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeReservedInstancesOfferingsOutput
		var err error
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			resp, err = conn.VM.DescribeReservedInstancesOfferings(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		})
		if err != nil {
			log.Printf("[DEBUG] Error reading lin (%s)", err)
		}
		if err != nil {
			return err
		}

		if len(resp.ReservedInstancesOfferings) > 0 && *resp.ReservedInstancesOfferings[0].ReservedInstancesOfferingId == rs.Primary.ID {
			return nil
		}

		return fmt.Errorf("Security Group not found")
	}
}

const testAccOutscaleOAPIReservedVmsOfferPurchaseEgressConfig = `
		resource "outscale_reserved_vms_offer_purchase" "test" {
			instance_count = 1
			reserved_vm_offering_id = ""
		}
	`
