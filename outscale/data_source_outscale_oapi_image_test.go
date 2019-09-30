package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIImageDataSource_Instance(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			skipIfNoOAPI(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIImageConfigBasic(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImageDataSourceID("data.outscale_image.nat_ami"),
					resource.TestCheckResourceAttr("data.outscale_image.nat_ami", "architecture", "x86_64"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIImageDataSource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIImageDataSourceBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIImageDataSourceID("data.outscale_image.omi"),
					testAccCheckState("data.outscale_image.omi"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIImageDataSourceID(n string) resource.TestCheckFunc {
	// Wait for IAM role
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find AMI data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AMI data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPIImageDataSourceBasicConfig = `
	data "outscale_image" "omi" {
		filter {
				name = "architectures"
				values = ["x86_64"]
		}
	}
`

func testAccCheckOutscaleOAPIImageConfigBasic(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = "${outscale_net.outscale_net.net_id}"
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
		}

		resource "outscale_vm" "basic" {
			image_id			           = "%s"
			vm_type                  = "%s"
			keypair_name	           = "terraform-basic"
			placement_subregion_name = "%sa"
			subnet_id                = "${outscale_subnet.outscale_subnet.subnet_id}"
			private_ips              =  ["10.0.0.12"]
		}

		resource "outscale_image" "foo" {
			image_name = "myImageName"
			vm_id = "${outscale_vm.basic.id}"
		}

		data "outscale_image" "nat_ami" {
			image_id = "${outscale_image.foo.id}"
		}
	`, omi, vmType, region)
}
