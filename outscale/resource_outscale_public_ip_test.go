package outscale

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOthers_PublicIP_basic(t *testing.T) {
	t.Parallel()
	var conf oscgo.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "outscale_public_ip.bar",
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		CheckDestroy:             testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},
		},
	})
}

func TestAccVM_PublicIP_instance(t *testing.T) {
	var conf oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	//rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "outscale_public_ip.bar1",
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		CheckDestroy:             testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPInstanceConfig(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar1", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},
		},
	})
}

// // This test is an expansion of TestAccOutscalePublicIP_instance, by testing the
// // associated Private PublicIPs of two instances
func TestAccNet_PublicIP_associated_user_private_ip(t *testing.T) {
	var one oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		IDRefreshName:            "outscale_public_ip.bar",
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		CheckDestroy:             testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPInstanceConfigAssociated(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},
			{
				Config: testAccOutscalePublicIPInstanceConfigAssociatedSwitch(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},
		},
	})
}

func testAccCheckOutscalePublicIPDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip" {
			continue
		}
		//Missing on Swagger Spec
		// if strings.Contains(rs.Primary.ID, "reservation") {
		// 	req := oscgo.ReadPublicIpsRequest{
		// 		Filters: oscgo.FiltersPublicIpcIp{
		// 			ReservationIds: []string{rs.Primary.ID},
		// 		},
		// 	}

		// 	var response *oscgo.ReadPublicIpsResponse
		// 	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		// 		var err error
		// 		resp, err := conn.oscgo.POST_ReadPublicIps(req)
		// 		response = resp.OK
		// 		return resource.RetryableError(err)
		// 	})

		// 	if err != nil {
		// 		// Verify the error is what we want
		// 		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
		// 			return nil
		// 		}

		// 		return err
		// 	}

		// 	if len(response.PublicIps) > 0 {
		// 		return fmt.Errorf("still exists")
		// 	}
		// } else {
		req := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				PublicIpIds: &[]string{rs.Primary.ID},
			},
		}

		var response oscgo.ReadPublicIpsResponse
		var statusCode int
		err := resource.Retry(60*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			response = rp
			statusCode = httpResp.StatusCode
			return nil
		})

		if err != nil {
			// Verify the error is what we want
			if statusCode == http.StatusNotFound {
				return nil
			}
			return err
		}

		if len(response.GetPublicIps()) > 0 {
			return fmt.Errorf("still exists")
		}
		//}
	}

	return nil
}

func testAccCheckOutscalePublicIPAttributes(conf *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conf.GetPublicIp() == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckOutscalePublicIPExists(n string, res *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PublicIP ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient)

		//Missing on Swagger Spec
		// if strings.Contains(rs.Primary.ID, "link") {
		// 	req := oscgo.ReadPublicIpsRequest{
		// 		Filters: oscgo.FiltersPublicIp{
		// 			ReservationIds: []string{rs.Primary.ID},
		// 		},
		// 	}
		// 	response, err := conn.oscgo.POST_ReadPublicIps(req)

		// 	if err != nil {
		// 		return err
		// 	}

		// 	if len(response.OK.PublicIps) != 1 ||
		// 		response.OK.PublicIps[0].ReservationId != rs.Primary.ID {
		// 		return fmt.Errorf("PublicIP not found")
		// 	}
		// 	*res = response.OK.PublicIps[0]

		// } else {
		req := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				PublicIpIds: &[]string{rs.Primary.ID},
			},
		}

		var response oscgo.ReadPublicIpsResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			response, _, err = conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()

			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
				return nil
			}

			return err
		}

		if len(response.GetPublicIps()) != 1 ||
			response.GetPublicIps()[0].GetPublicIpId() != rs.Primary.ID {
			return fmt.Errorf("PublicIP not found")
		}
		*res = response.GetPublicIps()[0]
		//}

		return nil
	}
}

const testAccOutscalePublicIPConfig = `
resource "outscale_public_ip" "bar" {
	tags {
		key = "Name"
		value = "public_ip_test"
	}
}
`

func testAccOutscalePublicIPInstanceConfig(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_ip" {
			security_group_name = "sg_publicIp"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name       = "%[4]s"
			security_group_ids = [outscale_security_group.sg_ip.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_public_ip" "bar1" {}

		resource "outscale_public_ip_link" "public_ip_link" {
			vm_id     = outscale_vm.basic.vm_id
			public_ip = outscale_public_ip.bar1.public_ip
		}
	`, omi, vmType, region, keypair)
}

func testAccOutscalePublicIPInstanceConfigAssociated(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sgIP" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[4]s"
			security_group_ids = [outscale_security_group.sgIP.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_vm" "basic2" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[4]s"
			security_group_ids = [outscale_security_group.sgIP.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region, keypair)
}

func testAccOutscalePublicIPInstanceConfigAssociatedSwitch(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sgIP" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[4]s"
			security_group_ids = [outscale_security_group.sgIP.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_vm" "basic2" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "%[4]s"
			security_group_ids = [outscale_security_group.sgIP.security_group_id]
			placement_subregion_name = "%[3]sa"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region, keypair)
}
