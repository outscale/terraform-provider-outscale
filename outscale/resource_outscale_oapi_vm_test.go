package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPIVM_Basic(t *testing.T) {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	var server fcu.Instance

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleOAPIVMAttributes(&server),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "type", "t2.micro"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_Update(t *testing.T) {
	// var server fcu.Instance

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	var before fcu.Instance
	var after fcu.Instance

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &before),
					testAccCheckOutscaleOAPIVMAttributes(&before),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "type", "t2.micro"),
				),
			},
			{
				Config: testAccInstanceConfigUpdateOAPIVMKey(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVMExists("outscale_vm.basic", &after),
					testAccCheckOAPIVMNotRecreated(
						t, &before, &after),
				),
			},
		},
	})
}

func testAccCheckOAPIVMExists(n string, i *fcu.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOAPIVMExistsWithProviders(n string, i *fcu.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
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

			var resp *fcu.DescribeInstancesOutput
			var err error
			for {
				resp, err = client.FCU.VM.DescribeInstances(&fcu.DescribeInstancesInput{
					InstanceIds: []*string{aws.String(rs.Primary.ID)},
				})
				if err != nil {
					time.Sleep(10 * time.Second)
				} else {
					break
				}

			}

			if len(resp.Reservations) > 0 {
				*i = *resp.Reservations[0].Instances[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckOAPIVMNotRecreated(t *testing.T,
	before, after *fcu.Instance) resource.TestCheckFunc {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}
	return func(s *terraform.State) error {
		if *before.InstanceId != *after.InstanceId {
			t.Fatalf("Outscale VM IDs have changed. Before %s. After %s", *before.InstanceId, *after.InstanceId)
		}
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

		var resp *fcu.DescribeInstancesOutput
		var err error
		for {
			// Try to find the resource
			resp, err = conn.FCU.VM.DescribeInstances(&fcu.DescribeInstancesInput{
				InstanceIds: []*string{&rs.Primary.ID},
			})
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

func testAccCheckOutscaleOAPIVMExists(n string, i *fcu.Instance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleOAPIVMExistsWithProviders(n string, i *fcu.Instance, providers *[]*schema.Provider) resource.TestCheckFunc {
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

			if resp.Reservations == nil {
				return fmt.Errorf("Instance not found")
			}

			if len(resp.Reservations) > 0 {
				*i = *resp.Reservations[0].Instances[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckOutscaleOAPIVMAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *server.ImageId != "ami-8a6a0120" {
			return fmt.Errorf("Bad image_id: %s", *server.ImageId)
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIVMConfig_basic() string {
	return `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	type = "t2.micro"
}`
}

func testAccInstanceConfigUpdateOAPIVMKey() string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	type = "t2.micro"
	keypair_name = "TestKey"
}`)
}
