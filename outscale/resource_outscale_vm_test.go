package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVM_Basic(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm

	resourceName := "outscale_vm.basic"

	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := fmt.Sprintf("%sa", os.Getenv("OUTSCALE_REGION"))
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigBasic(omi, "tinav4.c2r2p2", region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists(resourceName, &server),
					testAccCheckVMAttributes(t, &server, omi),

					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", "tinav4.c2r2p2"),

					resource.TestCheckResourceAttr(resourceName, "nested_virtualization", "true"),
				),
			},
		},
	})
}

func TestAccVMBehavior_Basic(t *testing.T) {
	t.Parallel()
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMBehaviorConfigBasic(omi, "tinav4.c2r2p2", region, keypair, "high", "stop"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
			{
				Config: testAccCheckVMBehaviorConfigBasic(omi, "tinav4.c2r2p2", region, keypair, "highest", "restart"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_importBasic(t *testing.T) {
	var (
		server       oscgo.Vm
		resourceName = "outscale_vm.basic"
		omi          = os.Getenv("OUTSCALE_IMAGEID")
		region       = fmt.Sprintf("%sa", os.Getenv("OUTSCALE_REGION"))
		keypair      = os.Getenv("OUTSCALE_KEYPAIR")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigBasic(omi, "tinav4.c2r2p2", region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists(resourceName, &server),
					resource.TestCheckResourceAttr(resourceName, "image_id", omi),
					resource.TestCheckResourceAttr(resourceName, "vm_type", "tinav4.c2r2p2"),
					resource.TestCheckResourceAttr(resourceName, "keypair_name", keypair),
					resource.TestCheckResourceAttr(resourceName, "placement_subregion_name", region),
					resource.TestCheckResourceAttr(resourceName, "private_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckVMImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_ips", "request_id"},
			},
		},
	})
}

func testAccCheckVMImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func TestAccVM_withNicAttached(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigBasicWithNicAttached(omi, "tinav4.c2r2p2", region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_withTags(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateVMTags(omi, "tinav4.c2r2p2", region, "Terraform-VM", keypair, sgId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_withNics(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigBasicWithNics(omi, "tinav4.c2r2p2", keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_Update(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	var before oscgo.Vm
	var after oscgo.Vm

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateVMKey(omi, "tinav4.c2r2p2", region, keypair, sgId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &before),
					testAccCheckVMAttributes(t, &before, omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
			{
				Config: testAccVmsConfigUpdateVMKey(omi, "tinav4.c2r2p2", region, keypair, sgId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &after),
					testAccCheckVMNotRecreated(t, &before, &after),
					resource.TestCheckResourceAttr("outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_WithSubnet(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigWithSubnet(omi, "tinav4.c2r2p2", region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func TestAccVM_WithBlockDeviceMappings(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	vmType := "tinav4.c2r2p2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigWithBlockDeviceMappings(omi, vmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.basic", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", vmType),
				),
			},
		},
	})
}

func TestAccVM_DeletionProtectionUpdate(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDeletionProtectionUpdateBasic(omi, "true", keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "true"),
				),
			},
			{
				Config: testAccCheckDeletionProtectionUpdateBasic(omi, "false", keypair),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccVMTags_Update(t *testing.T) {
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	sgId := os.Getenv("OUTSCALE_SECURITYGROUPID")

	//TODO: check tags
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateVMTags(omi, "tinav4.c2r2p2", region, "Terraform-VM", keypair, sgId),
				//Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccVmsConfigUpdateVMTags(omi, "tinav4.c2r2p2", region, "Terraform-VM2", keypair, sgId),
				//Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccVM_WithNet(t *testing.T) {
	var server oscgo.Vm
	omi := os.Getenv("OUTSCALE_IMAGEID")
	region := os.Getenv("OUTSCALE_REGION")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")
	vmType := "tinav4.c2r2p2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMConfigWithNet(omi, vmType, region, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.outscale_vmnet", &server),
					testAccCheckVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vmnet", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vmnet", "vm_type", vmType),
				),
			},
		},
	})
}

func TestAccVM_multiBlockDeviceMapping(t *testing.T) {
	var server oscgo.Vm
	region := os.Getenv("OUTSCALE_REGION")
	omi := os.Getenv("OUTSCALE_IMAGEID")
	keypair := os.Getenv("OUTSCALE_KEYPAIR")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVMWithMultiBlockDeviceMapping(region, omi, keypair),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVMExists("outscale_vm.outscale_vm", &server),
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "image_id", omi),
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "vm_type", "tinav4.c2r2p2"),
				),
			},
		},
	})
}

func testAccCheckVMWithMultiBlockDeviceMapping(region, omi, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "example" {
			subregion_name = "%sa"
			size           = 1
		}

		resource "outscale_snapshot" "snapshot" {
			volume_id = outscale_volume.example.id
		}

		resource "outscale_vm" "outscale_vm" {
			image_id     = "%s"
			vm_type      = "tinav4.c2r2p2"
			keypair_name = "%s"

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
					volume_type           = "io1"
					iops                  = 150
					snapshot_id           = outscale_snapshot.snapshot.id
					delete_on_vm_deletion = false
				}
			}

			tags {
				key   = "name"
				value = "VM with multiple Block Device Mappings"
			}
		}
	`, region, omi, keypair)
}

func testAccCheckDeletionProtectionUpdateBasic(omi, deletionProtection, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "outscale_vm" {
			image_id            = "%[1]s"
			vm_type             = "tinav4.c2r2p2"
			keypair_name        = "%[3]s"
			deletion_protection = %[2]s
		}
	`, omi, deletionProtection, keypair)
}

//TODO: check if is needed
// func testAccCheckVMSecurityGroupsUpdated(t *testing.T, before, after *oscgo.Vm) resource.TestCheckFunc {
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

func testAccCheckVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckVMExistsWithProviders(n, i, &providers)
}

func getVMsFilterByVMID(vmID string) *oscgo.FiltersVm {
	return &oscgo.FiltersVm{
		VmIds: &[]string{vmID},
	}
}

func testAccCheckVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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

			client := provider.Meta().(*Client)

			var resp oscgo.ReadVmsResponse
			var err error
			err = resource.Retry(120*time.Second, func() *resource.RetryError {
				rp, httpResp, err := client.OSCAPI.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				}).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp.StatusCode, err)
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

func testAccCheckVMNotRecreated(t *testing.T, before, after *oscgo.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertNotEqual(t, before.VmId, after.VmId, "Outscale VM IDs have changed.")
		return nil
	}
}

func testAccCheckVMDestroy(s *terraform.State) error {
	return testAccCheckVMDestroyWithProvider(s, testAccProvider)
}

//TODO: check if it is needed
// func testAccCheckVMDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		for _, provider := range *providers {
// 			if provider.Meta() == nil {
// 				continue
// 			}
// 			if err := testAccCheckVMDestroyWithProvider(s, provider); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}
// }

func testAccCheckVMDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*Client)

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
				return utils.CheckThrottling(httpResp.StatusCode, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})

		if err != nil {
			return err
		}
		if err == nil {
			for _, i := range resp.GetVms() {
				if i.GetState() != "" && i.GetState() != "terminated" {
					return fmt.Errorf("Found unterminated instance: %s", i.GetVmId())
				}
			}
		}

		// Verify the error is what we want
		if err != nil && statusCode == utils.ResourceNotFound {
			continue
		}
		return err
	}
	return nil
}

func testAccCheckVMAttributes(t *testing.T, server *oscgo.Vm, omi string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertEqual(t, omi, server.GetImageId(), "Bad image_id.")
		return nil
	}
}

func testAccCheckVMConfigBasic(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key   = "Name"
				value = "testacc-vm-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = outscale_net.outscale_net.net_id
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			placement_subregion_name = "%[3]s"
			placement_tenancy        = "dedicated"
			nested_virtualization    = true
			subnet_id                = outscale_subnet.outscale_subnet.subnet_id
			private_ips              =  ["10.0.0.12"]

			tags {
				key   = "name"
				value = "Terraform-VM"
			}
		}`, omi, vmType, region, keypair)
}

func testAccCheckVMConfigBasicWithNicAttached(omi, vmType, region, keypair string) string {
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

		resource "outscale_security_group" "outscale_security_group8" {
			description         = "test vm with nic"
			security_group_name = "private-sg-1"
			net_id              = outscale_net.outscale_net.net_id
		}

		resource "outscale_nic" "outscale_nic5" {
			subnet_id = outscale_subnet.outscale_subnet.subnet_id
		}

		resource "outscale_vm" "basic" {
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
			}

			nics {
				nic_id = outscale_nic.outscale_nic5.nic_id
				device_number = 1
			}

		}`, omi, vmType, region, keypair)
}

func testAccCheckVMConfigBasicWithNics(omi, vmType, keypair string) string {
	return fmt.Sprintf(`resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
	  }

	  resource "outscale_subnet" "outscale_subnet" {
		net_id         = outscale_net.outscale_net.net_id
		ip_range       = "10.0.0.0/24"
		subregion_name = "eu-west-2a"
	  }

	  resource "outscale_nic" "outscale_nic" {
		subnet_id = outscale_subnet.outscale_subnet.subnet_id
	  }

	  resource "outscale_security_group" "outscale_security_group" {
		description         = "test vm with nic"
		security_group_name = "private-sg"
		net_id              = outscale_net.outscale_net.net_id
	  }

	  resource "outscale_vm" "basic" {
		image_id     = "%s"
		vm_type      = "%s"
		keypair_name = "%s"

		# subnet_id              = outscale_subnet.outscale_subnet.subnet_id
		nics {
		  # delete_on_vm_deletion      = false
		  # description                = "myDescription"
		  device_number = 0

		  # nic_id                     = outscale_nic.outscale_nic.nic_id
		  # secondary_private_ip_count = 1
		  subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"

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
	  }`, omi, vmType, keypair)
}

func testAccVmsConfigUpdateVMKey(omi, vmType, region, keypair, sgId string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[4]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = outscale_net.net.id
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[4]s"
			security_group_ids       = ["%[5]s"]
			placement_subregion_name = "%[3]sb"
		}
	`, omi, vmType, region, keypair, sgId)
}

func testAccVmsConfigUpdateVMTags(omi, vmType string, region, value, keypair, sgId string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "net" {
			ip_range = "10.0.0.0/16"

			tags {
				key = "Name"
				value = "testacc-security-group-rs"
			}
		}

		resource "outscale_security_group" "sg" {
			security_group_name = "%[5]s"
			description         = "Used in the terraform acceptance tests"

			tags {
				key   = "Name"
				value = "tf-acc-test"
			}

			net_id = outscale_net.net.id
		}

		resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "%[5]s"
			security_group_ids       = ["%[6]s"]
			placement_subregion_name = "%[3]sb"

			tags {
				key   = "name"
				value = "%[4]s"
			}
		}
	`, omi, vmType, region, value, keypair, sgId)
}

func testAccCheckVMConfigWithSubnet(omi, vmType, region, keypair string) string {
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

func testAccCheckVMConfigWithBlockDeviceMappings(omi, vmType, region, keypair string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "external1" {
		subregion_name = "%[3]sa"
		size           = 1
	  }

	  resource "outscale_snapshot" "snapshot" {
		volume_id = outscale_volume.external1.id
	  }

	  resource "outscale_vm" "basic" {
		image_id     = "%[1]s"
		vm_type      = "%[2]s"
		keypair_name = "%[4]s"

		block_device_mappings {
		  device_name = "/dev/sdb"
		  no_device   = "/dev/xvdb"
		  bsu {
			volume_size           = 15
			volume_type           = "standard"
			snapshot_id           = outscale_snapshot.snapshot.id
			delete_on_vm_deletion = true
		  }
		}

		block_device_mappings {
		  device_name = "/dev/sdc"
		  bsu {
			volume_size           = 22
			volume_type           = "io1"
			iops                  = 150
			snapshot_id           = outscale_snapshot.snapshot.id
			delete_on_vm_deletion = true
		  }
		}

		block_device_mappings {
		  device_name = "/dev/sdc"
		  bsu {
			volume_size = 22
			volume_type = "io1"
			iops        = 150
			snapshot_id = outscale_snapshot.snapshot.id
		  }
		}
	  }
	`, omi, vmType, region, keypair)
}

func testAccCheckVMConfigWithNet(omi, vmType, region, keypair string) string {
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
		subregion_name = "%[3]sb"

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
		net_id = "${outscale_net.outscale_net.net_id}"

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
		net_id              = "${outscale_net.outscale_net.net_id}"
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

func testAccCheckVMBehaviorConfigBasic(omi, vmType, region, keypair, perfomance, vmBehavior string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"

			tags {
				key   = "Name"
				value = "testacc-vm-rs"
			}
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = outscale_net.outscale_net.net_id
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
		}

		resource "outscale_vm" "basic" {
			image_id			           = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name	           = "%[4]s"
			#placement_subregion_name = "%[3]s"
			subnet_id                = outscale_subnet.outscale_subnet.subnet_id
			private_ips              =  ["10.0.0.12"]
			vm_initiated_shutdown_behavior = "%[6]s"

			performance	           = "%[5]s"
			tags {
				key   = "name"
				value = "Terraform-VM"
			}
		}`, omi, vmType, region, keypair, perfomance, vmBehavior)
}
