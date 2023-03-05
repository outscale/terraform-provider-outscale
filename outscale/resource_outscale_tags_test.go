package outscale

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIVM_tags(t *testing.T) {
	v := &oscgo.Vm{}
	omi := os.Getenv("OUTSCALE_IMAGEID")

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckOAPIInstanceConfigTags(omi, "tinav4.c2r2p2", utils.GetRegion(), "keyOriginal", "valueOriginal"),
					Check: resource.ComposeTestCheckFunc(
						oapiTestAccCheckOutscaleVMExists("outscale_vm.vm", v),
						testAccCheckOAPIVMTags(v, "keyOriginal", "valueOriginal"),
						// Guard against regression of https://github.com/hashicorp/terraform/issues/914
						resource.TestCheckResourceAttr(
							"outscale_tag.foo", "tags.#", "1"),
					),
				},
				{
					Config: testAccCheckOAPIInstanceConfigTags(omi, "tinav4.c2r2p2", utils.GetRegion(), "keyUpdated", "valueUpdated"),
					Check: resource.ComposeTestCheckFunc(
						oapiTestAccCheckOutscaleVMExists("outscale_vm.vm", v),
						testAccCheckOAPIVMTags(v, "keyUpdated", "valueUpdated"),
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

func testAccCheckOAPIVMTags(vm *oscgo.Vm, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tags := vm.GetTags()
		return checkOAPITags(tags, key, value)
	}
}

func oapiTestAccCheckOutscaleVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return oapiTestAccCheckOutscaleVMExistsWithProviders(n, i, &providers)
}

func oapiTestAccCheckOutscaleVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			conn := provider.Meta().(*OutscaleClient)
			var resp oscgo.ReadVmsResponse
			var err error

			err = resource.Retry(30*time.Second, func() *resource.RetryError {
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
				return fmt.Errorf("VM not found")
			}

			if len(resp.GetVms()) > 0 {
				*i = resp.GetVms()[0]
				log.Printf("[DEBUG] VMS READ %+v", i)
				return nil
			}
		}
		return fmt.Errorf("VM not found")
	}
}

func testAccCheckOAPITags(
	ts *[]oscgo.ResourceTag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expected := map[string]string{
			"key":   key,
			"value": value,
		}
		tags := tagsOSCAPIToMap(*ts)
		for _, tag := range tags {
			if diff := deep.Equal(tag, expected); diff != nil {
				continue
			}
			return nil
		}
		return fmt.Errorf("error checking tags expected tag %+v is not found in %+v", expected, tags)
	}
}

func checkOAPITags(ts []oscgo.ResourceTag, key, value string) error {
	m := tagsOSCAPIToMap(ts)
	log.Printf("[DEBUG], tagsOAPIToMap=%+v", m)
	tag := m[0]

	if tag["key"] != key || tag["value"] != value {
		return fmt.Errorf("bad value expected: map[key:%s value:%s] got %+v", key, value, tag)
	}
	return nil
}

func testAccCheckOAPIInstanceConfigTags(omi, vmType, region, key, value string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-tags-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = outscale_net.outscale_net.net_id
			ip_range            = "10.0.0.0/24"
			subregion_name      = "%[3]sa"
		}
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%[3]sa"
			subnet_id                = outscale_subnet.outscale_subnet.subnet_id
			private_ips              =  ["10.0.0.12"]
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
