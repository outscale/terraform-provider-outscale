package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("outscale_vm", &resource.Sweeper{
		Name: "outscale_vm",
		F:    testSweepServers,
	})
}

func testSweepServers(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*OutscaleClient)

	var vms *fcu.DescribeInstancesOutput
	for {
		vms, err = client.FCU.VM.DescribeInstances(nil)
		if err != nil {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	var instanceids []*string

	fmt.Println("Before terminating sleep!")
	time.Sleep(1 * time.Second)

	for _, r := range vms.Reservations {
		for _, i := range r.Instances {
			if strings.HasPrefix(*i.KeyName, "terraform-") {
				instanceids = append(instanceids, i.KeyName)
			}
		}
	}

	for {
		_, err := client.FCU.VM.TerminateInstances(&fcu.TerminateInstancesInput{InstanceIds: instanceids})
		if err != nil {
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	return nil
}

func TestAccOutscaleServer_Basic(t *testing.T) {
	var server fcu.Instance

	var before fcu.Instance
	var after fcu.Instance

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleServerConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "instance_name", fmt.Sprintf("terraform-%d", rInt)),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "instance_type", "t2.micro"),
				),
			},
			{
				Config: testAccInstanceConfigUpdateVMKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("outscale_vm.basic", &after),
					testAccCheckInstanceNotRecreated(
						t, &before, &after),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string, i *fcu.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckInstanceExistsWithProviders(n, i, &providers)
}

func testAccCheckInstanceExistsWithProviders(n string, i *fcu.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			resp, err := client.FCU.VM.DescribeInstances(&fcu.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(rs.Primary.ID)},
			})
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidInstanceID.NotFound" {
				continue
			}
			if err != nil {
				return err
			}

			if len(resp.Reservations) > 0 {
				*i = *resp.Reservations[0].Instances[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckInstanceNotRecreated(t *testing.T,
	before, after *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *before.InstanceId != *after.InstanceId {
			t.Fatalf("Outscale VM IDs have changed. Before %s. After %s", *before.InstanceId, *after.InstanceId)
		}
		return nil
	}
}

func testAccCheckOutscaleVMDestroy(s *terraform.State) error {
	return testAccCheckOutscaleVMDestroyWithProvider(s, testAccProvider)
}

func testAccCheckOutscaleVMDestroyWithProviders(providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, provider := range *providers {
			if provider.Meta() == nil {
				continue
			}
			if err := testAccCheckOutscaleVMDestroyWithProvider(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckOutscaleVMDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_vm" {
			continue
		}

		var resp *fcu.DescribeInstancesOutput
		var err error
		for {
			// Try to find the resource
			resp, err = conn.FCU.VM.DescribeInstances(&fcu.DescribeInstancesInput{
				InstanceIds: []*string{&rs.Primary.ID},
			})
			if err != nil {
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}

		if err == nil {
			for _, r := range resp.Reservations {
				for _, i := range r.Instances {
					if i.State != nil && *i.State.Name != "terminated" {
						return fmt.Errorf("Found unterminated instance: %s", *i.InstanceId)
					}
				}
			}
		}

		// Verify the error is what we want
		if ae, ok := err.(awserr.Error); ok && ae.Code() == "InvalidInstanceID.NotFound" {
			continue
		}
		return err
	}

	return nil
}

func testAccCheckOutscaleVMExists(n string, i *fcu.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleVMExistsWithProviders(n string, i *fcu.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			var resp *fcu.DescribeInstancesOutput
			var err error

			for {
				resp, err = conn.FCU.VM.DescribeInstances(&fcu.DescribeInstancesInput{
					InstanceIds: []*string{&rs.Primary.ID},
				})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}
			}

			if fcuErr, ok := err.(awserr.Error); ok && fcuErr.Code() == "InvalidInstanceID.NotFound" {
				continue
			}
			if err != nil {
				return err
			}

			if len(resp.Reservations) > 0 {
				*i = *resp.Reservations[0].Instances[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckOutscaleServerAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *server.ImageId != "ami-8a6a0120" {
			return fmt.Errorf("Bad image_id: %s", *server.ImageId)
		}

		return nil
	}
}

func testAccCheckOutscaleServerConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	instance_name = "terraform-%d"
}`, rInt)
}

const testAccInstanceConfigUpdateVMKey = `
resource "outscale_vm" "basic" {
	image_id = "ami-5ad76458"
	instance_type = "t2.micro"
	key_name = "TestKey"
}`
