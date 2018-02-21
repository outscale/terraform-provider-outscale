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
	"github.com/terraform-providers/terraform-provider-outscale/osc"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var instanceId *string

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

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var server fcu.Instance

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleServerConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "group_set.#", "1"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "instances_set.0.group_set.#", "1"),
				),
			},
		},
	})
}

func getCredentials() *fcu.Client {
	config := osc.Config{
		Credentials: &osc.Credentials{
			AccessKey: "7E4U4AQ0CGLTWB78Q38V",
			SecretKey: "TDKLDVCNFDWFT6CVYBM9OPQ5YO9ZAJBN0JBJS99K",
			Region:    "eu-west-2",
		},
	}

	c, err := fcu.NewFCUClient(config)
	if err != nil {
		fmt.Println(err)
	}

	return c
}

func TestAccOutscaleServer_StopVM(t *testing.T) {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	c := getCredentials()

	keyname := "TestKey"
	var maxC int64
	imageID := "ami-8a6a0120"
	maxC = 1
	instanceType := "t2.micro"
	input := fcu.RunInstancesInput{
		ImageId:      &imageID,
		MaxCount:     &maxC,
		MinCount:     &maxC,
		KeyName:      &keyname,
		InstanceType: &instanceType,
	}
	output, err := c.VM.RunInstance(&input)
	fmt.Println(err)
	fmt.Println(output)

	time.Sleep(120 * time.Second)

	instanceId = output.Instances[0].InstanceId

	output3, err := c.VM.StopInstances(&fcu.StopInstancesInput{
		InstanceIds: []*string{instanceId},
	})

	fmt.Println(output3)

	if err != nil {
		t.Fatal(err)
	}
}
func TestAccOutscaleServer_StartVM(t *testing.T) {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	c := getCredentials()

	output3, err := c.VM.StartInstances(&fcu.StartInstancesInput{
		InstanceIds: []*string{instanceId},
	})

	fmt.Println(output3)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccOutscaleServer_UpdateVM(t *testing.T) {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	c := getCredentials()

	output3, err := c.VM.ModifyInstanceAttribute(&fcu.ModifyInstanceAttributeInput{
		InstanceId: instanceId,
		DisableApiTermination: &fcu.AttributeBooleanValue{
			Value: aws.Bool(false),
		},
	})
	fmt.Println(output3)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAccOutscaleServer_Windows_Password(t *testing.T) {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var server fcu.Instance

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleServerConfig_basic_windows(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basic_windows", &server),
					testAccCheckOutscaleWindowsServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic_windows", "image_id", "ami-e1b93f29"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic_windows", "instance_type", "t2.micro"),
				),
			},
		},
	})
}

func TestAccOutscaleServer_Update(t *testing.T) {
	// var server fcu.Instance

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	var before fcu.Instance
	var after fcu.Instance

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleServerConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleVMExists("outscale_vm.basic", &before),
					testAccCheckOutscaleServerAttributes(&before),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-8a6a0120"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "instance_type", "t2.micro"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "key_name", "terraform-basic"),
				),
			},
			{
				Config: testAccInstanceConfigUpdateVMKey(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("outscale_vm.basic", &after),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "key_name", "terraform-update"),
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

func testAccCheckInstanceNotRecreated(t *testing.T,
	before, after *fcu.Instance) resource.TestCheckFunc {

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
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

func testAccCheckOutscaleServerAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *server.ImageId != "ami-8a6a0120" {
			return fmt.Errorf("Bad image_id: %s", *server.ImageId)
		}

		if server.IpAddress == nil {
			return fmt.Errorf("No IP address found")
		}

		if len(*server.IpAddress) == 0 {
			return fmt.Errorf("Empty IP Address")
		}

		return nil
	}
}

func testAccCheckOutscaleWindowsServerAttributes(server *fcu.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *server.ImageId != "ami-e1b93f29" {
			return fmt.Errorf("Bad image_id: %s", *server.ImageId)
		}

		if server.IpAddress == nil {
			return fmt.Errorf("No IP address found")
		}

		if len(*server.IpAddress) == 0 {
			return fmt.Errorf("Empty IP Address")
		}

		return nil
	}
}

func testAccCheckOutscaleServerConfig_basic() string {
	return `
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}`
}

func testAccCheckOutscaleServerConfig_basic_windows() string {
	return `
resource "outscale_vm" "basic_windows" {
	image_id = "ami-e1b93f29"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}`
}

func testAccInstanceConfigUpdateVMKey() string {
	return fmt.Sprintf(`
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	security_group = ["sg-6ed31f3e"]
	key_name = "terraform-update"
}`)
}
