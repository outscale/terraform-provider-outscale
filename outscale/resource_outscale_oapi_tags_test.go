package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
)

func TestAccOutscaleOAPIVM_tags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapiFlag, err := strconv.ParseBool(o)
	if err != nil {
		oapiFlag = false
	}

	if !oapiFlag {
		t.Skip()
	}
	var v oapi.Vm

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOAPIInstanceConfigTags,
				Check: resource.ComposeTestCheckFunc(
					oapiTestAccCheckOutscaleVMExists("outscale_vm.foo", &v),
					testAccCheckOAPITags(v.Tags, "foo", "bar"),
					// Guard against regression of https://github.com/hashicorp/terraform/issues/914
					testAccCheckOAPITags(v.Tags, "#", ""),
				),
			},
		},
	})
}

func oapiTestAccCheckOutscaleVMExists(n string, i *oapi.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return oapiTestAccCheckOutscaleVMExistsWithProviders(n, i, &providers)
}

func oapiTestAccCheckOutscaleVMExistsWithProviders(n string, i *oapi.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			var resp *oapi.POST_ReadVmsResponses
			var err error

			for {
				resp, err = conn.OAPI.POST_ReadVms(oapi.ReadVmsRequest{
					Filters: oapi.FiltersVm{
						VmIds: []string{rs.Primary.ID},
					},
				})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			}

			if fcuErr, ok := err.(awserr.Error); ok && fcuErr.Code() == "InvalidInstanceID.NotFound" {
				continue
			}
			if err != nil {
				return err
			}

			if resp.OK.Vms == nil {
				return fmt.Errorf("Instance not found")
			}

			if len(resp.OK.Vms) > 0 {
				*i = resp.OK.Vms[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckOAPITags(
	ts []oapi.ResourceTag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := tagsOAPIToMap(ts)
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

const testAccCheckOAPIInstanceConfigTags = `
resource "outscale_vm" "foo" {
	image_id = "ami-8a6a0120"
	type = "m1.small"
}

resource "outscale_tag" "foo" {
	resource_ids = ["${outscale_vm.foo.id}"]
	tag = {
		faz = "baz"
	}
}
`
