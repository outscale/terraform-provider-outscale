package outscale

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIImage_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")
	omi := getOMIByRegion("eu-west-2", "ubuntu").OMI

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}
	var ami oapi.Image
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIImageConfigBasic(omi, "c4.large", rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPIImageExists("outscale_image.foo", &ami),
					testAccCheckState("outscale_image.foo"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "image_name", fmt.Sprintf("tf-testing-%d", rInt)),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.device_name", "/dev/sda1"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "block_device_mappings.0.bsu.delete_on_vm_deletion", "true"),
					resource.TestCheckResourceAttr(
						"outscale_image.foo", "state_comment.state_code", ""),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIImageRegisterConfig_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIImageRegisterConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccOutscaleOAPIImageCopyConfig_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOapi, err := strconv.ParseBool(o)
	if err != nil {
		isOapi = false
	}

	if !isOapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIImageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOAPIImageCopyConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckOAPIImageDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_image" {
			continue
		}

		// Try to find the OMI
		log.Printf("OMI-ID: %s", rs.Primary.ID)
		DescribeAmiOpts := &oapi.ReadImagesRequest{
			Filters: oapi.FiltersImage{ImageIds: []string{rs.Primary.ID}},
		}

		var result *oapi.ReadImagesResponse
		var resp *oapi.POST_ReadImagesResponses
		var err error

		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadImages(*DescribeAmiOpts)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					fmt.Printf("[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("Error retrieving Outscale Images: %s", errString)
		}

		result = resp.OK

		if len(result.Images) > 0 {
			state := result.Images[0].State
			return fmt.Errorf("OMI %s still exists in the state: %s", result.Images[0].ImageId, state)
		}
	}
	return nil
}

func testAccCheckOAPIImageExists(n string, ami *oapi.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("OMI Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OMI ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		DescribeAmiOpts := &oapi.ReadImagesRequest{
			Filters: oapi.FiltersImage{ImageIds: []string{rs.Primary.ID}},
		}

		var result *oapi.ReadImagesResponse
		var resp *oapi.POST_ReadImagesResponses
		var err error

		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadImages(*DescribeAmiOpts)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					fmt.Printf("[INFO] Request limit exceeded")
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})

		var errString string

		if err != nil || resp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if resp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
			} else if resp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
			} else if resp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
			}

			return fmt.Errorf("Error retrieving Outscale Images: %s", errString)
		}

		result = resp.OK

		if len(result.Images) == 0 {
			return fmt.Errorf("OMI not found")
		}
		*ami = result.Images[0]
		return nil
	}
}

func testAccOAPIImageConfigBasic(omi, vmType string, rInt int) string {
	return fmt.Sprintf(`
		resource "outscale_vm" "basic" {
			image_id			      = "%s"
			vm_type             = "%s"
			keypair_name		    = "terraform-basic"
			#security_group_ids = ["sg-6ed31f3e"]
		}

		resource "outscale_image" "foo" {
			image_name  = "tf-testing-%d"
			vm_id       = "${outscale_vm.basic.id}"
			no_reboot   = "true"
			description = "terraform testing"
		}
	`, omi, vmType, rInt)
}

const testAccOAPIImageRegisterConfig = `
resource "outscale_image" "outscale_image_register"
{ description = "Terraform-register-OMI"
image_name = "terraform-OMI-register"
file_location ="http://osu.eu-west-2.outscale.com/new-export-omi/omi-for-terraform/ami-b7d7f165/manifest?AWSAccessKeyId=S6AZO8TT4DOY9GOUOQ3U&Expires=1571754478&Signature=eOvObi6%2BFDW0AHEqafg5hZsLeJ4%3D"
# root_device_name= "/dev/sda1" (should not be used for register image)
}
`

const testAccOAPIImageCopyConfig = `
resource "outscale_image" "outscale_image_copy" {
	description = "Terraform-copy-OMI"
	image_name = "terraform-OMI-copy"
	source_image_id= "ami-3aa9428e"
	source_region_name= "eu-west-2"
	}
	`
