package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleKeypairsDataSource_Instance(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleKeypairsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeypairsDataSourceID("data.outscale_keypairs.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_keypairs.nat_ami", "key_set.0.key_name", "TestKey"),
				),
			},
		},
	})
}

func testAccCheckOutscaleKeypairsDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find AMI data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Key Pair data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleKeypairsDataSourceConfig = `
data "outscale_keypairs" "nat_ami" {
	filter {
		name = "key-name"
		values = ["TestKey"]
	}
}
`
