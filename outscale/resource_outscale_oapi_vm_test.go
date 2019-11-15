package outscale

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/outscale/osc-go/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVM_Basic(t *testing.T) {
	var server oapi.Vm
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
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
	var server oapi.Vm
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
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
				Config: testAccVmsConfigUpdateOAPIVMTags(omi, "c4.large", region),
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
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
	omi2 := getOMIByRegion("eu-west-2", "centos").OMI
	region := os.Getenv("OUTSCALE_REGION")

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
	var server oapi.Vm
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
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
	var server oapi.Vm
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI
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
						"outscale_vm.basic", vmType, "t2.micro"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_DeletionProtectionUpdate(t *testing.T) {
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

func testAccCheckOAPIVMSecurityGroupsUpdated(t *testing.T, before, after *oapi.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] ATTRS: %+v, %+v", before.SecurityGroups, after.SecurityGroups)
		if len(after.SecurityGroups) > 0 && len(before.SecurityGroups) > 0 {
			expectedSecurityGroup := after.SecurityGroups[0].SecurityGroupId
			for i := range before.SecurityGroups {
				assertNotEqual(t, before.SecurityGroups[i].SecurityGroupId, expectedSecurityGroup,
					"Outscale VM SecurityGroupId Either not found or are the same.")
			}
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
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
		}

		resource "outscale_subnet" "outscale_subnet" {
			net_id              = "${outscale_net.outscale_net.net_id}"
			ip_range            = "10.0.0.0/24"
			subregion_name      = "eu-west-2a"
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
			image_id			           = "%s"
			vm_type                  = "%s"
			keypair_name		         = "terraform-basic"
			# subnet_id              ="${outscale_subnet.outscale_subnet.subnet_id}"
			nics = [
				{
					# delete_on_vm_deletion      = false
					# description                = "myDescription"
					device_number                =  0
					# nic_id                     = "${outscale_nic.outscale_nic.nic_id}"
					# secondary_private_ip_count = 1
					subnet_id                    = "${outscale_subnet.outscale_subnet.subnet_id}"
					security_group_ids           = ["${outscale_security_group.outscale_security_group.security_group_id}"]					
					subnet_id                    = "${outscale_subnet.outscale_subnet.subnet_id}"
				  private_ips                  = [ 
				  	{
				  		private_ip = "10.0.0.123"
				  		is_primary = true   
						},
						{
				  		private_ip = "10.0.0.124"
				  		is_primary = false   
				  	}
				  ]
				}
			]
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

func testAccVmsConfigUpdateOAPIVMTags(omi, vmType string, region string) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id                 = "%s"
			vm_type                  = "%s"
			keypair_name             = "terraform-basic"
			security_group_ids       = ["sg-f4b1c2f8"]
			placement_subregion_name = "%sb"

			tags = {
				key   = "name"
				value = "terraform-subnet"
			}
		}
	`, omi, vmType, region)
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
			count = 1
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
			size = 1
		}
		
		resource "outscale_snapshot" "snapshot" {
			volume_id = "${outscale_volume.external1.id}"
		}

		resource "outscale_vm" "basic" {
			image_id              = "%[1]s"
			vm_type               = "%[2]s"
			keypair_name          = "terraform-basic"
	    block_device_mappings = [
				{
					device_name = "/dev/sdb"
					no_device   =  "/dev/xvdb"
					bsu = {
						volume_size=15
						volume_type = "gp2"
						snapshot_id = "${outscale_snapshot.snapshot.id}"
						delete_on_vm_deletion = true
					}
				},
				{
					device_name = "/dev/sdc"
					bsu = {
						volume_size=22
						volume_type = "io1"
						iops      = 150
						snapshot_id = "${outscale_snapshot.snapshot.id}"
						delete_on_vm_deletion = false
					}
				},
				{
					device_name = "/dev/sdc"
					bsu = {
						volume_size=22
						volume_type = "io1"
						iops      = 150
						snapshot_id = "${outscale_snapshot.snapshot.id}"
					}
				}
			]
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
