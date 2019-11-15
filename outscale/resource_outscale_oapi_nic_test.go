package outscale

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIENI_basic(t *testing.T) {
	var conf oapi.Nic
	subregion := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPINICDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf, subregion),
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccOutscaleOAPIENIConfigUpdate(subregion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf, subregion),
					resource.TestCheckResourceAttr("outscale_nic.outscale_nic", "private_ips.#", "3"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIENIExists(n string, res *oapi.Nic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ENI ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		dnir := &oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{NicIds: []string{rs.Primary.ID}},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(*dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || describeResp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if describeResp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
			} else if describeResp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
			} else if describeResp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
			}
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		result := describeResp.OK
		if len(result.Nics) != 1 ||
			result.Nics[0].NicId != rs.Primary.ID {
			return fmt.Errorf("ENI not found")
		}

		*res = result.Nics[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIENIAttributes(conf *oapi.Nic, suregion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if !reflect.DeepEqual(conf.LinkNic, oapi.LinkNic{}) {
			return fmt.Errorf("expected attachment to be nil")
		}

		if conf.SubregionName != fmt.Sprintf("%sa", suregion) {
			return fmt.Errorf("expected subregion_name to be %sa, but was %s", suregion, conf.SubregionName)
		}

		return nil
	}
}

func testAccOutscaleOAPIENIConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
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
		
			private_ips = [
				{
					is_primary = true
					private_ip = "10.0.0.23"
				},
				{
					is_primary = false
					private_ip = "10.0.0.46"
				},
			]
		}
	`, subregion)
}

func testAccOutscaleOAPIENIConfigUpdate(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
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
		
			private_ips = [
				{
					is_primary = true
					private_ip = "10.0.0.23"
				},
				{
					is_primary = false
					private_ip = "10.0.0.46"
				},
				{
					is_primary = false
					private_ip = "10.0.0.69"
				},
			]
		}	 
	`, subregion)
}
