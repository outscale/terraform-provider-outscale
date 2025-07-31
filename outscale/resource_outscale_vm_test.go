package outscale

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVM_Basic(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm

	resourceName := "outscale_vm.basic"

	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	region := fmt.Sprintf("%sa", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigBasic(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists(resourceName, &server),
					testAccCheckOutscaleVMAttributes(t, &server, omi),

					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
					resource.TestCheckResourceAttr(resourceName, "nested_virtualization", "false"),
				),
			},
		},
	})
}

func TestAccVM_uefi(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm

	resourceName := "outscale_vm.uefi"

	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	region := fmt.Sprintf("%sa", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMUefi(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists(resourceName, &server),
					testAccCheckOutscaleVMAttributes(t, &server, omi),

					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "boot_mode", "uefi"),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
					resource.TestCheckResourceAttr(resourceName, "nested_virtualization", "false"),
				),
			},
		},
	})
}

func TestAccVM_Behavior_Basic(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	region := fmt.Sprintf("%sa", utils.GetRegion())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMBehaviorConfigBasic(omi, utils.TestAccVmType, region, keypair, "high", "stop"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basicr1", &server),
					testAccCheckOutscaleVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basicr1", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basicr1", "vm_type", utils.TestAccVmType),
				),
			},
			{
				Config: testAccCheckOutscaleVMBehaviorConfigBasic(omi, utils.TestAccVmType, region, keypair, "high", "restart"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basicr1", &server),
					testAccCheckOutscaleVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basicr1", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basicr1", "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccVM_importBasic(t *testing.T) {
	var (
		server       oscgo.Vm
		resourceName = "outscale_vm.basic_import"
		omi          = os.Getenv("OUTSCALE_IMAGEID")
		keypair      = os.Getenv("OUTSCALE_KEYPAIR")
		region       = fmt.Sprintf("%sa", utils.GetRegion())
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigImport(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists(resourceName, &server),
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
					resource.TestCheckResourceAttr(resourceName, "keypair_name", keypair),
					resource.TestCheckResourceAttr(resourceName, "placement_subregion_name", region),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckOutscaleVMImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_ips", "request_id"},
			},
		},
	})
}

func testAccCheckOutscaleVMImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func TestAccNet_VM_withNicAttached(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	resourceName := "outscale_vm.basicNicAt"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigBasicWithNicAttached(omi, utils.TestAccVmType, utils.GetRegion(), keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccVM_withTags(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	tagsValue := "test_tags1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, utils.TestAccVmType, utils.GetRegion(), tagsValue, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basic_tags", &server),
					testAccCheckOutscaleVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic_tags", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic_tags", "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccNet_VM_withNics(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	resourceName := "outscale_vm.basic_with_nic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigBasicWithNics(omi, utils.TestAccVmType, keypair, utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccVM_UpdateKeypair(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	region := utils.GetRegion()
	resourceName := "outscale_vm.basic"

	var before oscgo.Vm
	var after oscgo.Vm

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMKey(omi, utils.TestAccVmType, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists(resourceName, &before),
					testAccCheckOutscaleVMAttributes(t, &before, omi),
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
				),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMKey2(omi, utils.TestAccVmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVMExists(resourceName, &after),
					testAccCheckOAPIVMNotRecreated(t, &before, &after),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccNet_VM_WithSubnet(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	resourceName := "outscale_vm.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigWithSubnet(omi, utils.TestAccVmType, utils.GetRegion(), keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func TestAccVM_UpdateDeletionProtection(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	vmType := utils.TestAccVmType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, keypair, vmType, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm1", "deletion_protection", "true"),
				),
			},
			{
				Config: testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, keypair, vmType, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm1", "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccVM_UpdateTags(t *testing.T) {
	t.Parallel()
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	tagsValue := "test_tags1"

	//TODO: check tags
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, utils.TestAccVmType, utils.GetRegion(), tagsValue, keypair),
				//Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, utils.TestAccVmType, utils.GetRegion(), "Terraform-VM2", keypair),
				//Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccNet_WithVM_PublicIp_Link(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"
	vmType := utils.TestAccVmType
	resourceName := "outscale_vm.outscale_vmnet"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMConfigWithNet(omi, vmType, utils.GetRegion(), keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", vmType),
				),
			},
		},
	})
}

func TestAccVM_multiBlockDeviceMapping(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := "terraform-basic"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: defineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleVMWithMultiBlockDeviceMapping(utils.GetRegion(), omi, keypair, utils.TestAccVmType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.outscale_vm", &server),
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "image_id", omi),
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "vm_type", utils.TestAccVmType),
				),
			},
		},
	})
}

func testAccCheckOutscaleVMWithMultiBlockDeviceMapping(region, omi, keypair, vmType string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%[1]sa"
			size           = 1
		}

		resource "outscale_snapshot" "snapshot" {
			volume_id = outscale_volume.example.id
		}

		resource "outscale_security_group" "sg_device" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sgProtection"
		}

		resource "outscale_vm" "outscale_vm" {
			image_id     = "%[2]s"
			vm_type      = "%[4]s"
			keypair_name = "%[3]s"
			security_group_ids = [outscale_security_group.sg_device.security_group_id]

			block_device_mappings {
				device_name = "/dev/sda1" # resizing bootdisk volume
				bsu {
					volume_size           = 100
					volume_type           = "standard"
					delete_on_vm_deletion = true
				}
			}

			block_device_mappings {
				device_name = "/dev/sdb"
				bsu {
					volume_size           = 30
					volume_type           = "gp2"
					snapshot_id           = outscale_snapshot.snapshot.id
					delete_on_vm_deletion = false
					tags {
						key           = "name"
						value         = "bsu-tags-gp2"
					}
				}
			}

			block_device_mappings {
				device_name = "/dev/sdc"
				bsu {
					volume_size = 100
					volume_type = "gp2"
					snapshot_id = outscale_snapshot.snapshot.id
					delete_on_vm_deletion = true
				}
			}

			tags {
				key   = "name"
				value = "VM with multiple Block Device Mappings"
			}
		}
	`, region, omi, keypair, vmType)
}

func testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, keypair, vmType string, deletionProtection bool) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_protection" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sgProtection"
		}

		resource "outscale_vm" "outscale_vm1" {
			image_id            = "%[1]s"
			vm_type             = "%[3]s"
			keypair_name        = "%[2]s"
			deletion_protection = %[4]t
			security_group_ids = [outscale_security_group.sg_protection.security_group_id]

		}
	`, omi, keypair, vmType, deletionProtection)
}

//TODO: check if is needed
// func testAccCheckOAPIVMSecurityGroupsUpdated(t *testing.T, before, after *oscgo.Vm) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		log.Printf("[DEBUG] ATTRS: %+v, %+v", before.GetSecurityGroups(), after.GetSecurityGroups())
// 		if len(after.GetSecurityGroups()) > 0 && len(before.GetSecurityGroups()) > 0 {
// 			expectedSecurityGroup := after.GetSecurityGroups()[0].GetSecurityGroupId()
// 			for i := range before.GetSecurityGroups() {
// 				assertNotEqual(t, before.GetSecurityGroups()[i].GetSecurityGroupId(), expectedSecurityGroup,
// 					"Outscale VM SecurityGroupId Either not found or are the same.")
// 			}
// 		}
// 		return nil
// 	}
// }

func testAccCheckOAPIVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOAPIVMExistsWithProviders(n, i, &providers)
}

func getVMsFilterByVMID(vmID string) *oscgo.FiltersVm {
	return &oscgo.FiltersVm{
		VmIds: &[]string{vmID},
	}
}

func testAccCheckOAPIVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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

			var resp oscgo.ReadVmsResponse
			var err error
			err = resource.Retry(120*time.Second, func() *resource.RetryError {
				rp, httpResp, err := client.OSCAPI.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				}).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				resp = rp
				return nil
			})

			if err != nil {
				return err
			}

			if len(resp.GetVms()) > 0 {
				*i = resp.GetVms()[0]
				return nil
			}
		}

		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOAPIVMNotRecreated(t *testing.T, before, after *oscgo.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertNotEqual(t, before.VmId, after.VmId, "Outscale VM IDs have changed.")
		return nil
	}
}

func testAccCheckOutscaleVMDestroy(s *terraform.State) error {
	return testAccCheckOutscaleVMDestroyWithProvider(s, testAccProvider)
}

//TODO: check if it is needed
// func testAccCheckOutscaleVMDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		for _, provider := range *providers {
// 			if provider.Meta() == nil {
// 				continue
// 			}
// 			if err := testAccCheckOutscaleVMDestroyWithProvider(s, provider); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// }

func testAccCheckOutscaleVMDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vm" {
			continue
		}

		var resp oscgo.ReadVmsResponse
		var err error
		var statusCode int
		// Try to find the resource
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			rp, httpResp, err := conn.OSCAPI.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: getVMsFilterByVMID(rs.Primary.ID),
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})

		if err != nil {
			return err
		}
		for _, i := range resp.GetVms() {
			if i.GetState() != "" && i.GetState() != "terminated" {
				return fmt.Errorf("Found running instance: %s", i.GetVmId())
			}
		}

		// Verify the error is what we want
		if err != nil && statusCode == http.StatusNotFound {
			continue
		}
		return err
	}
	return nil
}

func testAccCheckOutscaleVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			var resp oscgo.ReadVmsResponse
			var err error

			err = resource.Retry(120*time.Second, func() *resource.RetryError {
				rp, httpResp, err := conn.OSCAPI.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				}).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				resp = rp
				return nil
			})

			if err != nil {
				return err
			}

			if resp.Vms == nil {
				return fmt.Errorf("Vms not found")
			}

			if len(resp.GetVms()) > 0 {
				*i = resp.GetVms()[0]
				return nil
			}
		}
		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOutscaleVMAttributes(t *testing.T, server *oscgo.Vm, omi string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertEqual(t, omi, server.GetImageId(), "Bad image_id.")
		return nil
	}
}

func testAccCheckOutscaleVMUefi(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_vm_uefi" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sg_vm_uefi"
		}

		resource "outscale_vm" "uefi" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			placement_subregion_name = "%[3]s"
			security_group_ids = [outscale_security_group.sg_vm_uefi.security_group_id]
			boot_mode = "uefi"
			tags {
				key   = "name"
				value = "terraform_vm_uefi"
			}
		}`, omi, vmType, region, keypair)
}

func testAccCheckOutscaleVMConfigBasic(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_basicVvm" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sgVm"
		}
		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			placement_subregion_name = "%[3]s"
			security_group_ids = [outscale_security_group.sg_basicVvm.security_group_id]
			tags {
				key   = "name"
				value = "Terraform-VM"
			}
		}`, omi, vmType, region, keypair)
}

func testAccCheckOutscaleVMConfigImport(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_import_vm" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sg_importVm"
		}
		resource "outscale_vm" "basic_import" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	         = "%[4]s"
			placement_subregion_name = "%[3]s"
			security_group_ids = [outscale_security_group.sg_import_vm.security_group_id]

			tags {
				key   = "name"
				value = "Terraform-VM-import"
			}
		}`, omi, vmType, region, keypair)
}

func testAccCheckOutscaleVMConfigBasicWithNicAttached(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-vm-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = outscale_net.outscale_net.net_id
			ip_range            = "10.0.0.0/24"
			subregion_name      = "%[3]sa"
		}
		resource "outscale_security_group" "security_group7" {
			description         = "test vm with nic"
			security_group_name = "sg_nic5"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_security_group" "outscale_security_group8" {
			description         = "test vm with nic"
			security_group_name = "private-sg-1"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic5" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
			security_group_ids = [outscale_security_group.security_group7.security_group_id]
		}

		resource "outscale_vm" "basicNicAt" {
			image_id	         = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	         = "%[4]s"

			nics {
				subnet_id = outscale_subnet.outscale_subnet.subnet_id
				security_group_ids = [outscale_security_group.outscale_security_group8.security_group_id]
				private_ips  {
					  private_ip ="10.0.0.123"
					  is_primary = true
				 }
				device_number = 0
				delete_on_vm_deletion = true
			}

			nics {
				nic_id = outscale_nic.outscale_nic5.nic_id
				device_number = 1
			}

		}`, omi, vmType, region, keypair)
}

func testAccCheckOutscaleVMConfigBasicWithNics(omi, vmType, keypair, region string) string {
	return fmt.Sprintf(`
	  resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
	  }

	  resource "outscale_subnet" "outscale_subnet" {
		net_id         = outscale_net.outscale_net.net_id
		ip_range       = "10.0.0.0/24"
		subregion_name = "%[4]sa"
	  }

	  resource "outscale_security_group" "outscale_security_group" {
		description         = "test vm with nic"
		security_group_name = "private-sg"
		net_id              = outscale_net.outscale_net.net_id
	  }

	  resource "outscale_vm" "basic_with_nic" {
		image_id     = "%[1]s"
		vm_type      = "%[2]s"
		keypair_name = "%[3]s"

		primary_nic {
		  device_number = 0
		  subnet_id = outscale_subnet.outscale_subnet.subnet_id
                  delete_on_vm_deletion = true
		  security_group_ids = [outscale_security_group.outscale_security_group.security_group_id]

		  private_ips {
			private_ip = "10.0.0.123"
			is_primary = true
		  }

		  private_ips {
			private_ip = "10.0.0.124"
			is_primary = false
		  }
		}
	  }`, omi, vmType, keypair, region)
}

func testAccVmsConfigUpdateOAPIVMKey(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "keypair01" {
		keypair_name = "terraform-keypair-create"
		}
		resource "outscale_security_group" "sg_keypair" {
			security_group_name = "sg_keypair"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = outscale_keypair.keypair01.keypair_name
			security_group_ids       = [outscale_security_group.sg_keypair.security_group_id]
			placement_subregion_name = "%[3]sb"
		}
	`, omi, vmType, region)
}

func testAccVmsConfigUpdateOAPIVMKey2(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_keypair" "keypair01" {
			keypair_name = "terraform-keypair-create"
		}

		resource "outscale_security_group" "sg_keypair" {
			security_group_name = "sg_keypair"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = [outscale_security_group.sg_keypair.security_group_id]
			placement_subregion_name = "%[3]sb"
		}
	`, omi, vmType, region, keypair)
}

func testAccVmsConfigUpdateOAPIVMTags(omi, vmType, region, value, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_tags_vm" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sgTagsVm"
		}
		resource "outscale_vm" "basic_tags" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[5]s"
			placement_subregion_name = "%[3]sb"
			security_group_ids = [outscale_security_group.sg_tags_vm.security_group_id]

			tags {
				key   = "name"
				value = "%[4]s"
			}
		}
	`, omi, vmType, region, value, keypair)
}

func testAccCheckOutscaleVMConfigWithSubnet(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-vm-rs"
			}
		}

	  resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[3]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = outscale_net.outscale_net.net_id
	  }

	  resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg1-test-group_test-net"
			net_id              = outscale_net.outscale_net.net_id
	  }

	  resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = [outscale_security_group.outscale_security_group.security_group_id]
			subnet_id                = outscale_subnet.outscale_subnet.subnet_id
			placement_subregion_name = "%[3]sa"
			placement_tenancy        = "default"
	  }
	`, omi, vmType, region, keypair)
}

func testAccCheckOutscaleVMConfigWithNet(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"

		tags  {
			key   = "name"
			value = "Terraform_net"
		}
	}
	resource "outscale_subnet" "outscale_subnet" {
		net_id         = outscale_net.outscale_net.net_id
		ip_range       = "10.0.0.0/24"
		subregion_name = "%[3]sa"

		tags {
			key   = "name"
			value = "Terraform_subnet"
		}
	}

	resource "outscale_security_group" "outscale_sg" {
		description         = "sg for terraform tests"
		security_group_name = "terraform-sg"
		net_id              = outscale_net.outscale_net.net_id
	}

	resource "outscale_internet_service" "outscale_internet_service" {
                depends_on = [outscale_net.outscale_net]
        }

	resource "outscale_route_table" "outscale_route_table" {
		net_id = outscale_net.outscale_net.net_id
		tags {
			key   = "name"
			value = "Terraform_RT"
		}
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		route_table_id  = outscale_route_table.outscale_route_table.route_table_id
		subnet_id       = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		internet_service_id = outscale_internet_service.outscale_internet_service.internet_service_id
		net_id              = outscale_net.outscale_net.net_id
	}

	resource "outscale_route" "outscale_route" {
		gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
		destination_ip_range = "0.0.0.0/0"
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}
	resource "outscale_vm" "outscale_vmnet" {
		image_id           = "%[1]s"
		vm_type            = "%[2]s"
		keypair_name       = "%[4]s"
		security_group_ids = [outscale_security_group.outscale_sg.security_group_id]
		subnet_id          = outscale_subnet.outscale_subnet.subnet_id
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_public_ip_link" "outscale_public_ip_link" {
		vm_id     = outscale_vm.outscale_vmnet.vm_id
		public_ip = outscale_public_ip.outscale_public_ip.public_ip
	}
	`, omi, vmType, region, keypair)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		t.Fatalf(message+" Expected: %s and %s to differ.", a, b)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		t.Fatalf(message+"Expected: %s, actual: %s", a, b)
	}
}

func testAccCheckOutscaleVMBehaviorConfigBasic(omi, vmType, region, keypair, perfomance, vmBehavior string) string {
	return fmt.Sprintf(`
		resource "outscale_security_group" "sg_behavior_vm" {
			description                  = "testAcc Terraform security group"
			security_group_name          = "sg_behaviorVm"
		}
		resource "outscale_vm" "basicr1" {
			image_id                       = "%[1]s"
			vm_type                        = "%[2]s"
			keypair_name	               = "%[4]s"
			placement_subregion_name       = "%[3]s"
			vm_initiated_shutdown_behavior = "%[6]s"
			performance	               = "%[5]s"
			security_group_ids = [outscale_security_group.sg_behavior_vm.security_group_id]

			tags {
				key   = "name"
				value = "Terraform-VM"
			}
		}`, omi, vmType, region, keypair, perfomance, vmBehavior)
}
