package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceOutscaleOAPIVMType_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceOutscaleOAPIVMType(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.outscale_vm_type.test_by_id", "name", "t2.micro"),
				),
			},
		},
	})
}

func testAccDataSourceOutscaleOAPIVMType(rInt int) string {
	return fmt.Sprintf(`
data "outscale_vm_type" "test_by_id" {
	filter {
		name = "name"
		values = ["t2.micro"]
	}
}
`)
}
