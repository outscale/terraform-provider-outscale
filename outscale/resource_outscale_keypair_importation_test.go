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

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleKeyPairImportation_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}
	var conf fcu.KeyPairInfo

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleKeyPairImportationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleKeyPairConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleKeyPairImportationExists("outscale_keypair_importation.a_key_pair", &conf),
				),
			},
		},
	})
}

func testAccCheckOutscaleKeyPairImportationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_keypair_importation" {
			continue
		}

		// Try to find key pair
		var resp *fcu.DescribeKeyPairsOutput
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.FCU.VM.DescribeKeyPairs(&fcu.DescribeKeyPairsInput{
				KeyNames: []*string{aws.String(rs.Primary.ID)},
			})

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(err)
		})

		if resp == nil {
			return nil
		}

		if err == nil {
			if len(resp.KeyPairs) > 0 {
				return fmt.Errorf("still exist.")
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidKeyPair.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckOutscaleKeyPairImportationExists(n string, res *fcu.KeyPairInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Printf("[DEBUG]: Resources %s", s.RootModule().Resources)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KeyPair Importation id is set")
		}

		return nil
	}
}

func testAccOutscaleKeyPairImportationConfig(r int) string {
	return fmt.Sprintf(
		`
resource "outscale_keypair_importation" "a_key_pair" {
	key_name   = "tf-acc-key-pair-%d"
		public_key_material = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
}
`, r)
}
