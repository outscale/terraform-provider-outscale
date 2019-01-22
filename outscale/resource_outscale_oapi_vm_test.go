package outscale

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVM_Basic(t *testing.T) {
	var server oapi.Vm
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasic(omi, "c4.large"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleOAPIVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "c4.large"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_Update(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	var before oapi.Vm
	var after oapi.Vm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasic(omi, "t2.micro"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &before),
					testAccCheckOutscaleOAPIVMAttributes(t, &before, omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "vm_type", "t2.micro"),
				),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMKey(omi, "t2.micro"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVMExists("outscale_vm.basic", &after),
					testAccCheckOAPIVMNotRecreated(t, &before, &after),
					testAccCheckOAPIVMSecurityGroups(t, &before, &after),
				),
			},
		},
	})
}

func testAccCheckOAPIVMSecurityGroups(t *testing.T, before, after *oapi.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] ATTRS: %+v, %+v", before.SecurityGroups, after.SecurityGroups)
		expectedSecurityGroup := after.SecurityGroups[0].SecurityGroupId
		for i := range before.SecurityGroups {
			assertNotEqual(t, before.SecurityGroups[i].SecurityGroupId, expectedSecurityGroup,
				"Outscale VM SecurityGroupId Either not found or are the same.")
		}
		return nil
	}
}

func testAccCheckOAPIVMExists(n string, i *oapi.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOAPIVMExistsWithProviders(n string, i *oapi.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			client := provider.Meta().(*OutscaleClient)

			var resp *oapi.ReadVmsResponse
			var r *oapi.POST_ReadVmsResponses
			var err error
			for {
				r, err = client.OAPI.POST_ReadVms(oapi.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				})
				resp = r.OK
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}

			}

			if len(resp.Vms) > 0 {
				*i = resp.Vms[0]
				return nil
			}
		}

		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOAPIVMNotRecreated(t *testing.T, before, after *oapi.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertEqual(t, before.VmId, after.VmId, "Outscale VM IDs have changed.")
		return nil
	}
}

func testAccCheckOutscaleOAPIVMDestroy(s *terraform.State) error {
	return testAccCheckOutscaleOAPIVMDestroyWithProvider(s, testAccProvider)
}

func testAccCheckOutscaleOAPIVMDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckOutscaleOAPIVMDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckOutscaleOAPIVMDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vm" {
			continue
		}

		var resp *oapi.ReadVmsResponse
		var r *oapi.POST_ReadVmsResponses
		var err error
		for {
			// Try to find the resource
			r, err = conn.OAPI.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: getVMsFilterByVMID(rs.Primary.ID),
			})
			resp = r.OK
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			} else {
				break
			}
		}

		if err == nil {
			for _, i := range resp.Vms {
				if i.State != "" && i.State != "terminated" {
					return fmt.Errorf("Found unterminated instance: %s", i.VmId)
				}
			}
		}

		// Verify the error is what we want
		if ae, ok := err.(awserr.Error); ok && ae.Code() == "InvalidVmsID.NotFound" {
			continue
		}
		return err
	}

	return nil
}

func testAccCheckOutscaleOAPIVMExists(n string, i *oapi.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleOAPIVMExistsWithProviders(n string, i *oapi.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			conn := provider.Meta().(*OutscaleClient)
			var resp *oapi.ReadVmsResponse
			var r *oapi.POST_ReadVmsResponses
			var err error

			for {
				r, err = conn.OAPI.POST_ReadVms(oapi.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				})
				resp = r.OK
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			}

			if oapiErr, ok := err.(awserr.Error); ok && oapiErr.Code() == "InvalidVmsID.NotFound" {
				continue
			}
			if err != nil {
				return err
			}

			if resp.Vms == nil {
				return fmt.Errorf("Vms not found")
			}

			if len(resp.Vms) > 0 {
				*i = resp.Vms[0]
				return nil
			}
		}

		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOutscaleOAPIVMAttributes(t *testing.T, server *oapi.Vm, omi string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertEqual(t, omi, server.ImageId, "Bad image_id.")
		return nil
	}
}

func testAccCheckOutscaleOAPIVMConfigBasic(omi, vmType string) string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
	image_id			= "%s"
	vm_type            	= "%s"
	keypair_name		= "terraform-basic"
	security_group_ids	= ["sg-9752b7a6"]
}`, omi, vmType)
}

func testAccVmsConfigUpdateOAPIVMKey(omi, vmType string) string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
  image_id = "%s"
  vm_type = "%s"
  keypair_name = "integ_sut_keypair"
  security_group_ids = ["sg-22fda224"]
}`, omi, vmType)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		t.Fatalf(message+"Expected: %s and %s to differ.", a, b)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		t.Fatalf(message+"Expected: %s, actual: %s", a, b)
	}
}
