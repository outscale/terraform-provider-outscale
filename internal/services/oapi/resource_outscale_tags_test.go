package oapi_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_tags(t *testing.T) {
	v := &oscgo.Vm{}
	omi := os.Getenv("OUTSCALE_IMAGEID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testacc.PreCheck(t) },
		Providers:    testacc.SDKProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOAPIInstanceConfigTags(omi, testAccVmType, utils.GetRegion(), "keyOriginal", "valueOriginal"),
				Check: resource.ComposeTestCheckFunc(
					oapiTestAccCheckOutscaleVMExists("outscale_vm.vm", v),
					// Guard against regression of https://github.com/hashicorp/terraform/issues/914
					resource.TestCheckResourceAttr(
						"outscale_tag.foo", "tags.#", "1"),
				),
			},
			{
				Config: testAccCheckOAPIInstanceConfigTags(omi, testAccVmType, utils.GetRegion(), "keyUpdated", "valueUpdated"),
				Check: resource.ComposeTestCheckFunc(
					oapiTestAccCheckOutscaleVMExists("outscale_vm.vm", v),
					// Guard against regression of https://github.com/hashicorp/terraform/issues/914
					resource.TestCheckResourceAttr(
						"outscale_tag.foo", "tags.#", "1"),
				),
			},
		},
	})
}

func oapiTestAccCheckOutscaleVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		conn := testacc.SDKProvider.Meta().(*client.OutscaleClient)
		var resp oscgo.ReadVmsResponse
		err := retry.Retry(30*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.OSCAPI.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{rs.Primary.ID},
				},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return err
		}

		if len(resp.GetVms()) == 0 {
			return fmt.Errorf("vm not found")
		}

		*i = resp.GetVms()[0]
		log.Printf("[DEBUG] VMS READ %+v", i)
		return nil
	}
}

func testAccCheckOAPIInstanceConfigTags(omi, vmType, region, key, value string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_tag" "foo" {
			resource_ids = [outscale_vm.vm.id]

			tag {
				key   = "%s"
				value = "%s"
			}
		}
	`, omi, vmType, region, key, value)
}
