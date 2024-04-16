package outscale

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNet_WithNic_basic(t *testing.T) {
	var conf oscgo.Nic
	subregion := utils.GetRegion()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleNICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleENIConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleENIAttributes(&conf, subregion),
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "2"),
				),
			},
			{
				Config: testAccOutscaleENIConfigUpdate(subregion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleENIAttributes(&conf, subregion),
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "3"),
				),
			},
		},
	})
}

func testAccCheckOutscaleENIExists(n string, res *oscgo.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
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
			errString := err.Error()
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		if len(resp.GetNics()) != 1 ||
			resp.GetNics()[0].GetNicId() != rs.Primary.ID {
			return fmt.Errorf("ENI not found")
		}

		*res = resp.GetNics()[0]

		return nil
	}
}

func testAccCheckOutscaleENIAttributes(conf *oscgo.Nic, suregion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !reflect.DeepEqual(conf.GetLinkNic(), oscgo.LinkNic{}) {
			return fmt.Errorf("expected attachment to be nil")
		}

		if conf.GetSubregionName() != fmt.Sprintf("%sa", suregion) {
			return fmt.Errorf("expected subregion_name to be %sa, but was %s", suregion, conf.GetSubregionName())
		}

		return nil
	}
}

func testAccOutscaleENIConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name" 
				value = "testacc-nic-rs"
			}	
		}
		
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
		}
		
		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = outscale_net.outscale_net.net_id
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		
			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}

			private_ips	{
				is_primary = false
				private_ip = "10.0.0.46"
			}
		}
	`, subregion)
}

func testAccOutscaleENIConfigUpdate(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			tags {
				key = "Name" 
				value = "testacc-nic-rs"
			}	
		}
		
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
			ip_range       = "10.0.0.0/24"
			net_id         = outscale_net.outscale_net.net_id
		}
		
		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = outscale_net.outscale_net.net_id
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id          = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		
			private_ips {
				is_primary = true
				private_ip = "10.0.0.23"
			}
			
			private_ips {
				is_primary = false
				private_ip = "10.0.0.46"
			}
			
			private_ips {
				is_primary = false
				private_ip = "10.0.0.69"
			}
		}	 
	`, subregion)
}
