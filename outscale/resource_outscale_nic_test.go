package outscale

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOutscaleOAPIENI_basic(t *testing.T) {
	t.Parallel()
	subregion := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPIENIConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "2"),
				),
			},
			{
				Config: testAccOutscaleOAPIENIConfigUpdate(subregion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "3"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPINICDestroy(s *terraform.State) error {
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

func testAccOutscaleOAPIENIConfig(subregion string) string {
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
			net_id         = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id          = "${outscale_subnet.outscale_subnet.subnet_id}"
			security_group_ids = ["${outscale_security_group.outscale_sg.security_group_id}"]
		
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

func testAccOutscaleOAPIENIConfigUpdate(subregion string) string {
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
			net_id         = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_security_group" "outscale_sg" {
			description         = "sg for terraform tests"
			security_group_name = "terraform-sg"
			net_id              = "${outscale_net.outscale_net.net_id}"
		}
		
		resource "outscale_nic" "outscale_nic" {
			subnet_id          = "${outscale_subnet.outscale_subnet.subnet_id}"
			security_group_ids = ["${outscale_security_group.outscale_sg.security_group_id}"]
		
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

func testAccCheckOutscaleOAPIENIDestroy(s *terraform.State) error {
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
