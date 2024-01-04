package outscale

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccNet_withNicLink_Basic(t *testing.T) {
	var conf oscgo.Nic
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleNicLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNicLinkConfigBasic(rInt, omi, "tinav4.c2r2p2", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					resource.TestCheckResourceAttr(
						"outscale_nic_link.outscale_nic_link", "device_number", "1"),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_link.outscale_nic_link", "vm_id"),
					resource.TestCheckResourceAttrSet(
						"outscale_nic_link.outscale_nic_link", "nic_id"),
				),
			},
		},
	})
}

func TestAccNet_ImportNicLink_Basic(t *testing.T) {
	resourceName := "outscale_nic_link.outscale_nic_link"
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleNicLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleNicLinkConfigBasic(rInt, omi, "tinav4.c2r2p2", region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_nic_link.outscale_nic_link", "device_number", "1"),
					resource.TestCheckResourceAttrSet("outscale_nic_link.outscale_nic_link", "vm_id"),
					resource.TestCheckResourceAttrSet("outscale_nic_link.outscale_nic_link", "nic_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleNicLinkStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func testAccCheckOutscaleNicLinkStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		log.Printf("LOG_ : %#+v\n", rs.Primary.Attributes["nic_id"])
		return rs.Primary.Attributes["nic_id"], nil
	}
}

func testAccCheckOutscaleNicLinkDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic_link" {
			continue
		}

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

func testAccOutscaleNicLinkConfigBasic(sg int, omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key   = "Name"
				value = "testacc-nic-link"
			}
		}

		resource "outscale_security_group" "outscale_security_group" {
			security_group_name = "terraform_test_%d"
			description         = "Used in the terraform acceptance tests"
			net_id              = outscale_net.net.id

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "vm" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = [outscale_security_group.outscale_security_group.id]
			placement_subregion_name = "%[4]sa"
			subnet_id                = outscale_subnet.outscale_subnet.id
		}

		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[4]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.net.id

			depends_on = [outscale_security_group.outscale_security_group]
		}

		resource "outscale_nic" "outscale_nic" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
		}

		resource "outscale_nic_link" "outscale_nic_link" {
			device_number = 1
			vm_id         = outscale_vm.vm.id
			nic_id        = outscale_nic.outscale_nic.id
		}
	`, sg, omi, vmType, region)
}
