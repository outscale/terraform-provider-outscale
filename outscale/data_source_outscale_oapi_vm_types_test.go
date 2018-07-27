package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVMTypes_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIVMTypes(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vm_types.test_by_id", "instance_type_set.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMTypes(rInt int) string {
	return fmt.Sprintf(`
data "outscale_vm_types" "test_by_id" {
	filter {
		name = "name"
		values = ["t2.micro"]
	}
}
`)
}
