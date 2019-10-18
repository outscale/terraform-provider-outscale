package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleImage_importBasic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	resourceName := "outscale_image.foo"

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccImageConfigBasic(rInt),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"associate_public_ip_address", "user_data", "security_group"},
			},
		},
	})
}

func testAccImageConfigBasic(rInt int) string {
	return fmt.Sprintf(`
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

resource "outscale_image" "foo" {
	name = "tf-testing-%d"
	instance_id = "${outscale_vm.basic.id}"
}
	`, rInt, rInt)
}
