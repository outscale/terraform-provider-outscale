package outscale

import (
	"context"
	"fmt"
	"strings"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithNicDataSource_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleENIAttributes(&conf, utils.GetRegion()),
				),
			},
		},
	})
}

func TestAccNet_WithNicDataSource_basicFilter(t *testing.T) {
	t.Parallel()
	var conf oscgo.Nic

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleENIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIDataSourceConfigFilter(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleENIAttributes(&conf, utils.GetRegion()),
				),
			},
		},
	})
}

func testAccCheckOutscaleENIDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		var resp oscgo.ReadNicsResponse
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		req := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{NicIds: &[]string{rs.Primary.ID}},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return err
		}

		if len(resp.GetNics()) != 0 {
			return fmt.Errorf("Nic is not destroyed yet")
		}
	}
	return nil
}

func testAccCheckOutscaleNICDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		dnir := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{NicIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadNicsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnir).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				return nil
			}
			errString := err.Error()
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		if len(resp.GetNics()) > 0 {
			return fmt.Errorf("Nic with id %s is not destroyed yet", rs.Primary.ID)
		}
	}

	return nil
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

	data "outscale_nic" "outscale_nic" {
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

	data "outscale_nic" "outscale_nic" {
		filter {
			name = "nic_ids"
			values = [outscale_nic.outscale_nic.nic_id]
		}
	}
`, subregion)
}
