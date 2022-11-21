package outscale

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVM_tags(t *testing.T) {
	v := &oscgo.Vm{}
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

	if os.Getenv("TEST_QUOTA") == "true" {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckVMDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckInstanceConfigTags(omi, "tinav4.c2r2p2", region, "keyOriginal", "valueOriginal"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVMExists("outscale_vm.vm", v),
						testAccCheckVMTags(v, "keyOriginal", "valueOriginal"),
						// Guard against regression of https://github.com/hashicorp/terraform/issues/914
						resource.TestCheckResourceAttr(
							"outscale_tag.foo", "tags.#", "1"),
					),
				},
				{
					Config: testAccCheckInstanceConfigTags(omi, "tinav4.c2r2p2", region, "keyUpdated", "valueUpdated"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckVMExists("outscale_vm.vm", v),
						testAccCheckVMTags(v, "keyUpdated", "valueUpdated"),
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

func testAccCheckVMTags(vm *oscgo.Vm, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tags := vm.GetTags()
		return checkTags(tags, key, value)
	}
}

func testAccCheckTags(
	ts *[]oscgo.ResourceTag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expected := map[string]string{
			"key":   key,
			"value": value,
		}
		tags := tagsToMap(*ts)
		for _, tag := range tags {
			if diff := deep.Equal(tag, expected); diff != nil {
				continue
			}
			return nil
		}
		return fmt.Errorf("error checking tags expected tag %+v is not found in %+v", expected, tags)
	}
}

func checkTags(ts []oscgo.ResourceTag, key, value string) error {
	m := tagsToMap(ts)
	log.Printf("[DEBUG], tagsToMap=%+v", m)
	tag := m[0]

	if tag["key"] != key || tag["value"] != value {
		return fmt.Errorf("bad value expected: map[key:%s value:%s] got %+v", key, value, tag)
	}
	return nil
}

func testAccCheckInstanceConfigTags(omi, vmType, region, key, value string) string {
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
			subregion_name      = "eu-west-2a"
		}
		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			placement_subregion_name = "%sa"
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
