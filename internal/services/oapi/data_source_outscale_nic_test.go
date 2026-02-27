package oapi_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func TestAccNet_WithNicDataSource_basic(t *testing.T) {
	resourceName := "data.outscale_nic.data_nic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName:            "outscale_nic.outscale_nic",
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfig(utils.GetRegion(), sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "available"),
				),
			},
		},
	})
}

func TestAccNet_WithNicDataSource_basicFilter(t *testing.T) {
	resourceName := "data.outscale_nic.data_nic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		IDRefreshName:            "outscale_nic.outscale_nic",
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfigFilter(utils.GetRegion(), sgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

func testAccOutscaleENIDataSourceConfig(subregion, sgName string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-nic-ds"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		subregion_name = "%sa"
		ip_range       = "10.0.0.0/24"
		net_id         = outscale_net.outscale_net.id
	}
	resource "outscale_security_group" "sgdatNic" {
		security_group_name = "%[2]s"
		description         = "Used in the terraform acceptance tests"
		tags {
			key   = "Name"
			value = "tf-acc-test"
		}
		net_id       = outscale_net.outscale_net.net_id
	}

	resource "outscale_nic" "outscale_nic" {
		subnet_id = outscale_subnet.outscale_subnet.id
		security_group_ids = [outscale_security_group.sgdatNic.security_group_id]
		tags {
			value = "tf-value"
			key   = "tf-key"
		}
	}

	data "outscale_nic" "data_nic" {
		filter {
			name = "nic_ids"
			values = [outscale_nic.outscale_nic.nic_id]
		}
	}
	`, subregion, sgName)
}

func testAccOutscaleENIDataSourceConfigFilter(subregion, sgName string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags {
			key = "Name"
			value = "testacc-nic-ds-filter"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		subregion_name = "%sa"
		ip_range       = "10.0.0.0/16"
		net_id         = outscale_net.outscale_net.id
	}
	resource "outscale_security_group" "sgdatNic" {
		security_group_name = "%[2]s"
		description         = "Used in the terraform acceptance tests"
		tags {
			key   = "Name"
			value = "tf-acc-test"
		}
		 net_id       = outscale_net.outscale_net.net_id
	}

	resource "outscale_nic" "outscale_nic" {
		subnet_id = outscale_subnet.outscale_subnet.id
		security_group_ids = [outscale_security_group.sgdatNic.security_group_id]
		tags {
			value = "tf-value"
			key   = "tf-key"
		}
	}

	data "outscale_nic" "data_nic" {
		filter {
			name = "nic_ids"
			values = [outscale_nic.outscale_nic.nic_id]
		}
	}
`, subregion, sgName)
}
