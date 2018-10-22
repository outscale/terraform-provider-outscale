package outscale

import (
	"fmt"
	"os"
	"strconv"
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
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	var server oapi.Vms_2

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &server),
					testAccCheckOutscaleOAPIVMAttributes(&server),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-b1d1f100"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "type", "t2.micro"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIVM_Update(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	var before oapi.Vms_2
	var after oapi.Vms_2

	// rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPIVMConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIVMExists("outscale_vm.basic", &before),
					testAccCheckOutscaleOAPIVMAttributes(&before),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "image_id", "ami-b1d1f100"),
					resource.TestCheckResourceAttr(
						"outscale_vm.basic", "type", "t2.micro"),
				),
			},
			{
				Config: testAccVmsConfigUpdateOAPIVMKey(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIVMExists("outscale_vm.basic", &after),
					testAccCheckOAPIVMNotRecreated(
						t, &before, &after),
				),
			},
		},
	})
}

func testAccCheckOAPIVMExists(n string, i *oapi.Vms_2) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOAPIVMExistsWithProviders(n string, i *oapi.Vms_2, providers *[]*schema.Provider) resource.TestCheckFunc {
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
					Filters: getVMsFiltersByVMID(rs.Primary.ID),
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

func testAccCheckOAPIVMNotRecreated(t *testing.T,
	before, after *oapi.Vms_2) resource.TestCheckFunc {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}
	return func(s *terraform.State) error {
		if before.VmId != after.VmId {
			t.Fatalf("Outscale VM IDs have changed. Before %s. After %s", before.VmId, after.VmId)
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

		var resp *oapi.ReadVmsResponse
		var r *oapi.POST_ReadVmsResponses
		var err error
		for {
			// Try to find the resource
			r, err = conn.OAPI.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: getVMsFiltersByVMID(rs.Primary.ID),
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

func testAccCheckOutscaleOAPIVMExists(n string, i *oapi.Vms_2) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOutscaleOAPIVMExistsWithProviders(n, i, &providers)
}

func testAccCheckOutscaleOAPIVMExistsWithProviders(n string, i *oapi.Vms_2, providers *[]*schema.Provider) resource.TestCheckFunc {
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
					Filters: getVMsFiltersByVMID(rs.Primary.ID),
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

func testAccCheckOutscaleOAPIVMAttributes(server *oapi.Vms_2) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if server.ImageId != "ami-b1d1f100" {
			return fmt.Errorf("Bad image_id: %s", server.ImageId)
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIVMConfigBasic() string {
	return `
resource "outscale_vm" "basic" {
	image_id = "ami-b1d1f100"
	type = "t2.micro"
}`
}

func testAccVmsConfigUpdateOAPIVMKey() string {
	return fmt.Sprintf(`
resource "outscale_vm" "outscale_vm" {
  image_id = "ami-b1d1f100"
  type = "c4.large"
  keypair_name = "integ_sut_keypair"
  #firewall_rules_set_ids = ["sg-c73d3b6b"]
  #firewall_rules_set_id = "sg-c73d3b6b" # tempo tests
}`)
}
