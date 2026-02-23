package oapi_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccVM_tags(t *testing.T) {
	v := &osc.Vm{}
	omi := os.Getenv("OUTSCALE_IMAGEID")

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.ParallelTest(t, resource.TestCase{
			Providers:    testacc.SDKProviders,
			CheckDestroy: testAccCheckOutscaleVMDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckOAPIInstanceConfigTags(omi, testAccVmType, utils.GetRegion(), "keyOriginal", "valueOriginal"),
					Check: resource.ComposeTestCheckFunc(
						oapiTestAccCheckOutscaleVMExists(t.Context(), "outscale_vm.vm", v),
						// Guard against regression of https://github.com/hashicorp/terraform/issues/914
						resource.TestCheckResourceAttr(
							"outscale_tag.foo", "tags.#", "1"),
					),
				},
				{
					Config: testAccCheckOAPIInstanceConfigTags(omi, testAccVmType, utils.GetRegion(), "keyUpdated", "valueUpdated"),
					Check: resource.ComposeTestCheckFunc(
						oapiTestAccCheckOutscaleVMExists(t.Context(), "outscale_vm.vm", v),
						// Guard against regression of https://github.com/hashicorp/terraform/issues/914
						resource.TestCheckResourceAttr(
							"outscale_tag.foo", "tags.#", "1"),
					),
				},
			},
		})
	} else {
		t.Skip("will be done soon")
	}
}

func oapiTestAccCheckOutscaleVMExists(ctx context.Context, n string, i *osc.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		client := testacc.SDKProvider.Meta().(*client.OutscaleClient).OSC
		resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
			Filters: &osc.FiltersVm{
				VmIds: &[]string{rs.Primary.ID},
			},
		}, options.WithRetryTimeout(DefaultTimeout))
		if err != nil {
			return err
		}

		if resp.Vms == nil || len(*resp.Vms) == 0 {
			return fmt.Errorf("vm not found")
		}

		*i = (*resp.Vms)[0]
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
