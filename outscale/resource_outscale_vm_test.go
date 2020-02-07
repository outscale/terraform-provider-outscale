package outscale

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVM_Basic(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasic(omi, "c4.large", region),
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

func TestAccOutscaleOAPIVM_BasicTags(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, "c4.large", region, "Terraform-VM"),
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

func TestAccOutscaleOAPIVM_BasicWithNics(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasicWithNics(omi, "c4.large"),
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
	region := os.Getenv("OUTSCALE_REGION")
	omi := getOMIByRegion(region, "centos").OMI
	omi2 := getOMIByRegion(region, "centos").OMI

	var before oscgo.Vm
	var after oscgo.Vm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMKey(omi, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &before),
					testAccCheckOutscaleOAPIVMAttributes(t, &before, omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr("outscale_vm.basic", "vm_type", "c4.large"),
				),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMKey(omi2, "c4.large", region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVMExists("outscale_vm.basic", &after),
					testAccCheckOAPIVMNotRecreated(t, &before, &after),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_WithSubnet(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigWithSubnet(omi, "c4.large", region),
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

func TestAccOutscaleOAPIVM_WithBlockDeviceMappings(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")
	vmType := "t2.micro"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigWithBlockDeviceMappings(omi, vmType, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleOAPIVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "vm_type", vmType),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_DeletionProtectionUpdate(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "centos").OMI

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "true"),
				),
			},
			{
				Config: testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_vm.outscale_vm", "deletion_protection", "false"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVMTags_Update(t *testing.T) {
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	//omi2 := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")

	//var before oscgo.Vm
	//var after oscgo.Vm

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, "c4.large", region, "Terraform-VM"),
				//Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, "c4.large", region, "Terraform-VM2"),
				//Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_WithNet(t *testing.T) {
	var server oscgo.Vm
	omi := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")
	vmType := "t2.micro"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigWithNet(omi, vmType, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.outscale_vmnet", &server),
					testAccCheckOutscaleOAPIVMAttributes(t, &server, omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vmnet", "image_id", omi),
					resource.TestCheckResourceAttr(
						"outscale_vm.outscale_vmnet", "vm_type", vmType),
				),
			},
		},
	})
}

func testAccCheckOutscaleDeletionProtectionUpdateBasic(omi, deletionProtection string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "outscale_vm" {
			image_id            = "%s"
			vm_type             = "c4.large"
			keypair_name        = "terraform-basic"
			deletion_protection = %s
		}
	`, omi, deletionProtection)
}

func testAccCheckOAPIVMSecurityGroupsUpdated(t *testing.T, before, after *oscgo.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] ATTRS: %+v, %+v", before.GetSecurityGroups(), after.GetSecurityGroups())
		if len(after.GetSecurityGroups()) > 0 && len(before.GetSecurityGroups()) > 0 {
			expectedSecurityGroup := after.GetSecurityGroups()[0].GetSecurityGroupId()
			for i := range before.GetSecurityGroups() {
				assertNotEqual(t, before.GetSecurityGroups()[i].GetSecurityGroupId(), expectedSecurityGroup,
					"Outscale VM SecurityGroupId Either not found or are the same.")
			}
		}
		return nil
	}
}

func testAccCheckOAPIVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOSCAPIVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOSCAPIVMExistsWithProviders(n, i, &providers)
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
			for {
				resp, _, err = client.OSCAPI.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				})})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}

			}

			if len(resp.GetVms()) > 0 {
				*i = resp.GetVms()[0]
				return nil
			}
		}

		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOSCAPIVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			for {
				resp, _, err = client.OSCAPI.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
					Filters: getOSCVMsFilterByVMID(rs.Primary.ID),
				})})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}

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

		var resp oscgo.ReadVmsResponse
		var err error
		for {
			// Try to find the resource
			resp, _, err = conn.OSCAPI.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
				Filters: getVMsFilterByVMID(rs.Primary.ID),
			})})
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
			for _, i := range resp.GetVms() {
				if i.GetState() != "" && i.GetState() != "terminated" {
					return fmt.Errorf("Found unterminated instance: %s", i.GetVmId())
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

func testAccCheckOutscaleOAPIVMExists(n string, i *oscgo.Vm) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleOAPIVMExistsWithProviders(n string, i *oscgo.Vm, providers *[]*schema.Provider) resource.TestCheckFunc {
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

			for {
				resp, _, err = conn.OSCAPI.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(oscgo.ReadVmsRequest{
					Filters: getVMsFilterByVMID(rs.Primary.ID),
				})})
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

			if len(resp.GetVms()) > 0 {
				*i = resp.GetVms()[0]
				return nil
			}
		}

		return fmt.Errorf("Vms not found")
	}
}

func testAccCheckOutscaleOAPIVMAttributes(t *testing.T, server *oscgo.Vm, omi string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assertEqual(t, omi, *server.ImageId, "Bad image_id.")
		return nil
	}
}

func testAccCheckOutscaleOAPIVMConfigBasic(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = "${outscale_net.outscale_net.net_id}"
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
		}

		resource "outscale_vm" "basic" {
			image_id			           = "%s"
			vm_type                  = "%s"
			keypair_name	           = "terraform-basic"
			placement_subregion_name = "%sa"
			subnet_id                = "${outscale_subnet.outscale_subnet.subnet_id}"
			private_ips              =  ["10.0.0.12"]
		}`, omi, vmType, region)
}

func testAccCheckOutscaleOAPIVMConfigBasicWithNics(omi, vmType string) string {
	return fmt.Sprintf(`resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
	  }
	  
	  resource "outscale_subnet" "outscale_subnet" {
		net_id         = "${outscale_net.outscale_net.net_id}"
		ip_range       = "10.0.0.0/24"
		subregion_name = "eu-west-2a"
	  }
	  
	  resource "outscale_nic" "outscale_nic" {
		subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
	  }
	  
	  resource "outscale_security_group" "outscale_security_group" {
		description         = "test vm with nic"
		security_group_name = "private-sg"
		net_id              = "${outscale_net.outscale_net.net_id}"
	  }
	  
	  resource "outscale_vm" "basic" {
		image_id     = "%s"
		vm_type      = "%s"
		keypair_name = "terraform-basic"
	  
		# subnet_id              ="${outscale_subnet.outscale_subnet.subnet_id}"
		nics {
		  # delete_on_vm_deletion      = false
		  # description                = "myDescription"
		  device_number = 0
	  
		  # nic_id                     = "${outscale_nic.outscale_nic.nic_id}"
		  # secondary_private_ip_count = 1
		  subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
	  
		  security_group_ids = ["${outscale_security_group.outscale_security_group.security_group_id}"]
	  
		  private_ips {
			private_ip = "10.0.0.123"
			is_primary = true
		  }
	  
		  private_ips {
			private_ip = "10.0.0.124"
			is_primary = false
		  }
		}
	  }`, omi, vmType)
}

func testAccVmsConfigUpdateOAPIVMKey(omi, vmType, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"
		}
	`, omi, vmType, region)
}

func testAccVmsConfigUpdateOAPIVMTags(omi, vmType string, region, value string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"

			tags {
				key   = "name"
				value = "%s"
			}
		}
	`, omi, vmType, region, value)
}

func testAccCheckOutscaleOAPIVMConfigWithSubnet(omi, vmType string, region string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}
	  
	  resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%[3]sa"
			ip_range       = "10.0.0.0/16"
			net_id         = "${outscale_net.outscale_net.net_id}"
	  }
	  
	  resource "outscale_security_group" "outscale_security_group" {
			description         = "test group"
			security_group_name = "sg1-test-group_test-net"
			net_id              = "${outscale_net.outscale_net.net_id}"
	  }
	  
	  resource "outscale_vm" "basic" {
			image_id                 = "%[1]s"
			vm_type                  = "%[2]s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["${outscale_security_group.outscale_security_group.security_group_id}"]
			subnet_id                = "${outscale_subnet.outscale_subnet.subnet_id}"
			placement_subregion_name = "%sa"
			placement_tenancy        = "default"
	  }	  
	`, omi, vmType, region)
}

func testAccCheckOutscaleOAPIVMConfigWithBlockDeviceMappings(omi, vmType, region string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "external1" {
		subregion_name = "eu-west-2a"
		size           = 1
	  }
	  
	  resource "outscale_snapshot" "snapshot" {
		volume_id = "${outscale_volume.external1.id}"
	  }
	  
	  resource "outscale_vm" "basic" {
		image_id     = "%[1]s"
		vm_type      = "%[2]s"
		keypair_name = "terraform-basic"
	  
		block_device_mappings {
		  device_name = "/dev/sdb"
		  no_device   = "/dev/xvdb"
		  bsu = {
			volume_size           = 15
			volume_type           = "gp2"
			snapshot_id           = "${outscale_snapshot.snapshot.id}"
			delete_on_vm_deletion = true
		  }
		}
	  
		block_device_mappings {
		  device_name = "/dev/sdc"
		  bsu = {
			volume_size           = 22
			volume_type           = "io1"
			iops                  = 150
			snapshot_id           = "${outscale_snapshot.snapshot.id}"
			delete_on_vm_deletion = true
		  }
		}
	  
		block_device_mappings {
		  device_name = "/dev/sdc"
		  bsu = {
			volume_size = 22
			volume_type = "io1"
			iops        = 150
			snapshot_id = "${outscale_snapshot.snapshot.id}"
		  }
		}
	  }
	`, omi, vmType, region)
}

func testAccCheckOutscaleOAPIVMConfigWithNet(omi, vmType, region string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		
		tags  {                               
			key   = "name"
			value = "Terraform_net"
		}
	}
	resource "outscale_subnet" "outscale_subnet" {
		net_id         = "${outscale_net.outscale_net.net_id}"
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
		net_id              = "${outscale_net.outscale_net.net_id}"
	}
	
	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = "${outscale_net.outscale_net.net_id}"
		
		tags {                               
			key   = "name"
			value = "Terraform_RT"
		}
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		route_table_id  = "${outscale_route_table.outscale_route_table.route_table_id}"
		subnet_id       = "${outscale_subnet.outscale_subnet.subnet_id}"
	}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		internet_service_id = "${outscale_internet_service.outscale_internet_service.internet_service_id}" 
		net_id              = "${outscale_net.outscale_net.net_id}"
	}

	resource "outscale_route" "outscale_route" {
		gateway_id           = "${outscale_internet_service.outscale_internet_service.internet_service_id}"
		destination_ip_range = "0.0.0.0/0"
		route_table_id       = "${outscale_route_table.outscale_route_table.route_table_id}"
	} 
	resource "outscale_vm" "outscale_vmnet" {
		image_id           = "%[1]s"
		vm_type            = "%[2]s"
		keypair_name       = "terraform-basic"
		security_group_ids = ["${outscale_security_group.outscale_sg.security_group_id}"]
		subnet_id          ="${outscale_subnet.outscale_subnet.subnet_id}"
	}

	resource "outscale_public_ip" "outscale_public_ip" {}

	resource "outscale_public_ip_link" "outscale_public_ip_link" {
		vm_id     = "${outscale_vm.outscale_vmnet.vm_id}"
		public_ip = "${outscale_public_ip.outscale_public_ip.public_ip}"
	}
	`, omi, vmType, region)
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
