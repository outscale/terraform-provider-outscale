package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func TestAccNet_WithNicDataSource_basic(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_nic.data_nic"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "outscale_nic.outscale_nic",
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "available"),
				),
			},
		},
	})
}

func TestAccNet_WithNicDataSource_basicFilter(t *testing.T) {
	t.Parallel()
	resourceName := "data.outscale_nic.data_nic"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "outscale_nic.outscale_nic",
		ProtoV5ProviderFactories: defineTestProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfigFilter(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				),
			},
		},
	})
}

func testAccOutscaleENIDataSourceConfig(subregion string) string {
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
		security_group_name = "test_sgNic"
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
	`, subregion)
}

func testAccOutscaleENIDataSourceConfigFilter(subregion string) string {
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
		security_group_name = "test_sgNic"
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
`, subregion)
}
