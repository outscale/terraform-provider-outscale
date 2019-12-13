package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"log"
	"os"
	"testing"
	"time"

	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVM_tags(t *testing.T) {
	v := &oscgo.Vm{}
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOAPIInstanceConfigTags(omi, "c4.large", region, "keyOriginal", "valueOriginal"),
				Check: resource.ComposeTestCheckFunc(
					oapiTestAccCheckOutscaleVMExists("outscale_vm.vm", v),
					testAccCheckOAPIVMTags(v, "keyOriginal", "valueOriginal"),
					// Guard against regression of https://github.com/hashicorp/terraform/issues/914
					resource.TestCheckResourceAttr(
						"outscale_tag.foo", "tags.#", "1"),
				),
			},
			{
				Config: testAccCheckOAPIInstanceConfigTags(omi, "c4.large", region, "keyUpdated", "valueUpdated"),
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

			for {
				resp, _, err = conn.OSCAPI.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
					Filters: &oscgo.FiltersVm{
						VmIds: &[]string{rs.Primary.ID},
					},
				})})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			}

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
	ts []oscgo.ResourceTag, key string, value string) resource.TestCheckFunc {
	log.Printf("[DEBUG] testAccCheckOAPITags %+v", ts)
	return func(s *terraform.State) error {
		m := tagsOSCAPIToMap(ts)
		v, ok := m[0]["Key"]
		if value != "" && !ok {
			return fmt.Errorf("Missing tag: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra tag: %s", key)
		}
		if value == "" {
			return nil
		}
		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}
		return nil
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
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = "${outscale_net.outscale_net.net_id}"
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
		}
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%sa"
			subnet_id                = "${outscale_subnet.outscale_subnet.subnet_id}"
			private_ips              =  ["10.0.0.12"]
		}

		resource "outscale_tag" "foo" {
			resource_ids = ["${outscale_vm.vm.id}"]

			tag {
				key   = "%s"
				value = "%s"			
			}
		}
	`, omi, vmType, region, key, value)
}
