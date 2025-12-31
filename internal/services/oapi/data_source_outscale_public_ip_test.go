package oapi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccOthers_DataSourcePublicIP(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscalePublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceOutscalePublicIPCheck("data.outscale_public_ip.by_public_ip_id"),
					testAccDataSourceOutscalePublicIPCheck("data.outscale_public_ip.by_public_ip"),
				),
			},
		},
	})
}

func TestAccVM_WithPublicIP(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscalePublicIPConfigwithVM(omi, utils.GetRegion(), oapi.TestAccVmType, sgName),
			},
		},
	})
}

func testAccDataSourceOutscalePublicIPCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		eipRs, ok := s.RootModule().Resources["outscale_public_ip.test"]
		if !ok {
			return fmt.Errorf("can't find outscale_public_ip.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["public_ip_id"] != eipRs.Primary.Attributes["public_ip_id"] {
			return fmt.Errorf(
				"public_ip_id is %s; want %s",
				attr["public_ip_id"],
				eipRs.Primary.Attributes["public_ip_id"],
			)
		}

		if attr["public_ip"] != eipRs.Primary.Attributes["public_ip"] {
			return fmt.Errorf(
				"public_ip is %s; want %s",
				attr["public_ip"],
				eipRs.Primary.Attributes["public_ip"],
			)
		}

		return nil
	}
}

func TestAccOthers_DataSourcePublicIP_withTags(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutscalePublicIPConfigWithTags,
			},
		},
	})
}

const testAccDataSourceOutscalePublicIPConfig = `
	resource "outscale_public_ip" "test" {}

	data "outscale_public_ip" "by_public_ip_id" {
	  public_ip_id = outscale_public_ip.test.public_ip_id
	}

	data "outscale_public_ip" "by_public_ip" {
		filter {
			name = "public_ips"
			values = [outscale_public_ip.test.public_ip]
		}
	}
`

const testAccDataSourceOutscalePublicIPConfigWithTags = `
	resource "outscale_public_ip" "outscale_public_ip" {
		tags {
			key   = "name"
			value = "public_ip-data"
		}
	}

	data "outscale_public_ip" "outscale_public_ip" {
		filter {
			name   = "tags"
			values = ["name=public_ip-data"]
		}

		filter {
			name   = "public_ip_ids"
			values = [outscale_public_ip.outscale_public_ip.public_ip_id]
		}
	}
`

func testAccDataSourceOutscalePublicIPConfigwithVM(omi, region, vmType, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_Pbip" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "outscale_vm" {
			image_id     = "%[1]s"
			vm_type      = "%[3]s"
			keypair_name = "terraform-basic"
			security_group_ids = [outscale_security_group.sg_Pbip.security_group_id]
		}

		resource "outscale_public_ip" "outscale_public_ip" {
			tags {
				key   = "name"
				value = "Terraform_EIP"
			}
			tags {
				key   = "platform"
				value = "%[2]s"
			}
			tags {
				key   = "project"
				value = "terraform"
			}
		}

		resource "outscale_public_ip_link" "outscale_public_ip_link" {
			vm_id     = outscale_vm.outscale_vm.vm_id
			public_ip = outscale_public_ip.outscale_public_ip.public_ip
		}

		data "outscale_public_ip" "outscale_public_ip-5" {
			filter {
				name   = "link_public_ip_ids"
				values = [outscale_public_ip_link.outscale_public_ip_link.link_public_ip_id]
			}
		}
	`, omi, region, vmType, sgName)
}
