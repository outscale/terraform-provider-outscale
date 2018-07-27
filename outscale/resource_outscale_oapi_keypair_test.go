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

func TestAccOutscaleOAPIKeyPair_basic(t *testing.T) {
	var conf fcu.KeyPairInfo
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIKeyPairConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIKeyPairExists("outscale_keypair.a_key_pair", &conf),
					testAccCheckOutscaleOAPIKeyPairFingerprint("8a:47:95:bb:b1:45:66:ef:99:f5:80:91:cc:be:94:48", &conf),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIKeyPair_basic_name(t *testing.T) {
	var conf fcu.KeyPairInfo
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIKeyPairConfigRetrieveName(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIKeyPairExists("outscale_keypair.a_key_pair", &conf),
					resource.TestCheckResourceAttr(
						"outscale_keypair.a_key_pair", "key_name", "tf-acc-key-pair",
					),
				),
			},
		},
	})
}
func TestAccOutscaleOAPIKeyPair_generatedName(t *testing.T) {
	var conf fcu.KeyPairInfo

	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi == false {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIKeyPairConfigGeneratedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIKeyPairExists("outscale_keypair.a_key_pair", &conf),
					testAccCheckOutscaleOAPIKeyPairFingerprint("8a:47:95:bb:b1:45:66:ef:99:f5:80:91:cc:be:94:48", &conf),
					func(s *terraform.State) error {
						if conf.KeyName == nil {
							return fmt.Errorf("bad: No SG name")
						}
						if !strings.HasPrefix(*conf.KeyName, "terraform-") {
							return fmt.Errorf("No terraform- prefix: %s", *conf.KeyName)
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIKeyPairDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_keypair" {
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
				return fmt.Errorf("still exist")
			}
			return nil
		}

		// Verify the error is what we want
		ec2err, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if ec2err.Code() != "InvalidOAPIKeyPair.NotFound" {
			return err
		}
	}

	return nil
}

func testAccCheckOutscaleOAPIKeyPairFingerprint(expectedFingerprint string, conf *fcu.KeyPairInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.KeyFingerprint != expectedFingerprint {
			return fmt.Errorf("incorrect fingerprint. expected %s, got %s", expectedFingerprint, *conf.KeyFingerprint)
		}
		return nil
	}
}

func testAccCheckOutscaleOAPIKeyPairExists(n string, res *fcu.KeyPairInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OAPIKeyPair name is set")
		}
		var resp *fcu.DescribeKeyPairsOutput
		conn := testAccProvider.Meta().(*OutscaleClient)

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
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
		if len(resp.KeyPairs) != 1 ||
			*resp.KeyPairs[0].KeyName != rs.Primary.ID {
			return fmt.Errorf("OAPIKeyPair not found")
		}

		*res = *resp.KeyPairs[0]

		return nil
	}
}

func testAccCheckOutscaleOAPIKeyPairNamePrefix(t *testing.T) {
	var conf fcu.KeyPairInfo

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:        func() { testAccPreCheck(t) },
		IDRefreshName:   "outscale_keypair.a_key_pair",
		IDRefreshIgnore: []string{"key_name_prefix"},
		Providers:       testAccProviders,
		CheckDestroy:    testAccCheckOutscaleOAPIKeyPairDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckOutscaleOAPIKeyPairPrefixNameConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIKeyPairExists("outscale_keypair.a_key_pair", &conf),
					testAccCheckOutscaleOAPIKeyPairGeneratedNamePrefix(
						"outscale_keypair.a_key_pair", "baz-"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIKeyPairGeneratedNamePrefix(
	resource, prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		name, ok := r.Primary.Attributes["name"]
		if !ok {
			return fmt.Errorf("Name attr not found: %#v", r.Primary.Attributes)
		}
		if !strings.HasPrefix(name, prefix) {
			return fmt.Errorf("Name: %q, does not have prefix: %q", name, prefix)
		}
		return nil
	}
}

func testAccOutscaleOAPIKeyPairConfig(r int) string {
	return fmt.Sprintf(
		`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "tf-acc-key-pair-%d"
	key_material = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
}
`, r)
}

func testAccOutscaleOAPIKeyPairConfigRetrieveName(r int) string {
	return fmt.Sprintf(
		`
resource "outscale_keypair" "a_key_pair" {
	key_name   = "tf-acc-key-pair-%d"
}
`, r)
}

const testAccOutscaleOAPIKeyPairConfigGeneratedName = `
resource "outscale_keypair" "a_key_pair" {
	key_material = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
}
`

func testAccCheckOutscaleOAPIKeyPairPrefixNameConfig(r int) string {
	return fmt.Sprintf(
		`
resource "outscale_keypair" "a_key_pair" {
	key_name_prefix   = "baz-%d"
	key_material = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 phodgson@thoughtworks.com"
}
`, r)
}
