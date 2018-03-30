package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleVM_tags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var v fcu.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckInstanceConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.foo", &v),
					testAccCheckTags(&v.Tags, "foo", "bar"),
					testAccCheckTags(&v.Tags, "#", ""),
				),
			},
		},
	})
}

func testAccCheckTags(
	ts *[]*fcu.Tag, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		m := tagsToMap(*ts)
		k, ok := m[0]["key"]
		v, ok := m[0]["value"]
		if value != "" && !ok {
			return fmt.Errorf("Missing tag: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra tag: %s", key)
		}
		if value == "" {
			return nil
		}

		fmt.Println("KEY =>", key, "K =>", k)
		fmt.Println("V =>", v, "VALUE =>", value)

		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}

		return nil
	}
}

const testAccCheckInstanceConfigTags = `
resource "outscale_keypair" "a_key_pair" {
	key_name   = "terraform-key-%d"
}

resource "outscale_firewall_rules_set" "web" {
  group_name = "terraform_acceptance_test_example_1"
  group_description = "Used in the terraform acceptance tests"
}

resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	security_group = ["${outscale_firewall_rules_set.web.id}"]
	key_name = "${outscale_keypair.a_key_pair.key_name}"
}

resource "outscale_tag" "outscale_tag" {
    resource_ids = ["${outscale_vm.basic.id}"]

    tag {                              
        name7 = "testDataSource7"           
    }                                      
}

`
