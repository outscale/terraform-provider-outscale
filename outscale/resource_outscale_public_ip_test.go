package outscale

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPIPublicIP_basic(t *testing.T) {
	var conf oscgo.PublicIp

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIPublicIP_instance(t *testing.T) {
	var conf oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

	//rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfig(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfig2(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &conf),
					testAccCheckOutscaleOAPIPublicIPAttributes(&conf),
				),
			},
		},
	})
}

// // This test is an expansion of TestAccOutscalePublicIP_instance, by testing the
// // associated Private PublicIPs of two instances
func TestAccOutscaleOAPIPublicIP_associated_user_private_ip(t *testing.T) {
	var one oscgo.PublicIp
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_public_ip.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfigAssociated(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscaleOAPIPublicIPAttributes(&one),
				),
			},

			resource.TestStep{
				Config: testAccOutscaleOAPIPublicIPInstanceConfigAssociatedSwitch(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIPublicIPExists("outscale_public_ip.bar", &one),
					testAccCheckOutscaleOAPIPublicIPAttributes(&one),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIPublicIPDestroy(s *terraform.State) error {
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
				PublicIps: &[]string{rs.Primary.ID},
			},
		}

		var response oscgo.ReadPublicIpsResponse
		err := resource.Retry(60*time.Second, func() *resource.RetryError {
			var err error
			response, _, err = conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(req)})
			return resource.RetryableError(err)
		})

		if err != nil {
			// Verify the error is what we want
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
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

func testAccCheckOutscaleOAPIPublicIPAttributes(conf *oscgo.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conf.GetPublicIp() == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIPublicIPExists(n string, res *oscgo.PublicIp) resource.TestCheckFunc {
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
				PublicIps: &[]string{rs.Primary.ID},
			},
		}

		var response oscgo.ReadPublicIpsResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			response, _, err = conn.OSCAPI.PublicIpApi.ReadPublicIps(context.Background(), &oscgo.ReadPublicIpsOpts{ReadPublicIpsRequest: optional.NewInterface(req)})

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

		if err != nil {

			// Verify the error is what we want
			if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidPublicIps.NotFound") {
				return nil
			}

			return err
		}

		if len(response.GetPublicIps()) != 1 ||
			response.GetPublicIps()[0].GetPublicIp() != rs.Primary.ID {
			return fmt.Errorf("PublicIP not found")
		}
		*res = response.GetPublicIps()[0]
		//}

		return nil
	}
}

const testAccOutscaleOAPIPublicIPConfig = `
resource "outscale_public_ip" "bar" {
	tags {
		key = "Name"
		value = "public_ip_test"
	}
}
`

func testAccOutscaleOAPIPublicIPInstanceConfig(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region)
}

func testAccOutscaleOAPIPublicIPInstanceConfig2(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region)
}

func testAccOutscaleOAPIPublicIPInstanceConfigAssociated(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "terraform-basic"
			security_group_ids = ["sg-f4b1c2f8"]
			placement_subregion_name = "%[3]sb"
		}

		resource "outscale_vm" "basic2" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "terraform-basic"
			security_group_ids = ["sg-f4b1c2f8"]
			placement_subregion_name = "%[3]sb"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region)
}

func testAccOutscaleOAPIPublicIPInstanceConfigAssociatedSwitch(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "terraform-basic"
			security_group_ids = ["sg-f4b1c2f8"]
			placement_subregion_name = "%[3]sb"
		}

		resource "outscale_vm" "basic2" {
			image_id           = "%[1]s"
			vm_type            = "%[2]s"
			keypair_name       = "terraform-basic"
			security_group_ids = ["sg-f4b1c2f8"]
			placement_subregion_name = "%[3]sb"
		}

		resource "outscale_public_ip" "bar" {}
	`, omi, vmType, region)
}
