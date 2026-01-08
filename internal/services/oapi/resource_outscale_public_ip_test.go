package oapi_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_PublicIP_basic(t *testing.T) {
	var conf oscgo.PublicIp

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
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
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPInstanceConfig(omi, testAccVmType, region, keypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar1", &conf),
					testAccCheckOutscalePublicIPAttributes(&conf),
				),
			},
			// Ignore attributes related to the Public IP Link, that gets populated after a refresh
			testacc.ImportStep("outscale_public_ip.bar1", "link_public_ip_id", "nic_account_id", "nic_id", "private_ip", "vm_id", "request_id"),
		},
	})
}

// // This test is an expansion of TestAccOutscalePublicIP_instance, by testing the
// // associated Private PublicIPs of two instances
func TestAccNet_PublicIP_associated_user_private_ip(t *testing.T) {
	var one oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := utils.GetRegion()
	keypair := "terraform-basic"
	sgName := acctest.RandomWithPrefix("testacc-sg")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testacc.PreCheck(t) },
		IDRefreshName:            "outscale_public_ip.bar",
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckOutscalePublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscalePublicIPInstanceConfigAssociated(omi, testAccVmType, region, keypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},
			{
				Config: testAccOutscalePublicIPInstanceConfigAssociatedSwitch(omi, testAccVmType, region, keypair, sgName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscalePublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscalePublicIPAttributes(&one),
				),
			},
		},
	})
}

func testAccCheckOutscalePublicIPDestroy(s *terraform.State) error {
	client := testacc.ConfiguredClient.OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_public_ip" {
			continue
		}
		// Missing on Swagger Spec
		// if strings.Contains(rs.Primary.ID, "reservation") {
		// 	req := oscgo.ReadPublicIpsRequest{
		// 		Filters: oscgo.FiltersPublicIpcIp{
		// 			ReservationIds: []string{rs.Primary.ID},
		// 		},
		// 	}

		// 	var response *oscgo.ReadPublicIpsResponse
		// 	err := retry.Retry(60*time.Second, func() *retry.RetryError {
		// 		var err error
		// 		resp, err := conn.oscgo.POST_ReadPublicIps(req)
		// 		response = resp.OK
		// 		return retry.RetryableError(err)
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
		err := retry.Retry(60*time.Second, func() *retry.RetryError {
			rp, httpResp, err := client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
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
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no publicip id is set")
		}

		client := testacc.ConfiguredClient.OSCAPI

		// Missing on Swagger Spec
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
		// 		return fmt.Errorf("publicip not found")
		// 	}
		// 	*res = response.OK.PublicIps[0]

		// } else {
		req := oscgo.ReadPublicIpsRequest{
			Filters: &oscgo.FiltersPublicIp{
				PublicIpIds: &[]string{rs.Primary.ID},
			},
		}

		var response oscgo.ReadPublicIpsResponse
		err := retry.Retry(120*time.Second, func() *retry.RetryError {
			var err error
			response, _, err = client.PublicIpApi.ReadPublicIps(context.Background()).ReadPublicIpsRequest(req).Execute()
			if err != nil {
				if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
					return retry.RetryableError(err)
				}

				return retry.NonRetryableError(err)
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
			return fmt.Errorf("publicip not found")
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

func testAccOutscalePublicIPInstanceConfig(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_ip" {
			security_group_name = "%[5]s"
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
	`, omi, vmType, region, keypair, sgName)
}

func testAccOutscalePublicIPInstanceConfigAssociated(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sgIP" {
			security_group_name = "%[5]s"
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
	`, omi, vmType, region, keypair, sgName)
}

func testAccOutscalePublicIPInstanceConfigAssociatedSwitch(omi, vmType, region, keypair, sgName string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sgIP" {
			security_group_name = "%[5]s"
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
	`, omi, vmType, region, keypair, sgName)
}
