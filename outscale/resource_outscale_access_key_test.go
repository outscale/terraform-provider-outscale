package outscale

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestAccOutscaleAccessKey_basic(t *testing.T) {
	resourceName := "outscale_access_key.outscale_access_key"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleAccessKeyBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccOutscaleAccessKey_updatedToInactivedKey(t *testing.T) {
	resourceName := "outscale_access_key.outscale_access_key"

	state := "ACTIVE"
	stateUpdated := "INACTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleAccessKeyUpdatedConfig(state),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", state),
				),
			},
			{
				Config: testAccOutscaleAccessKeyUpdatedConfig(stateUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", stateUpdated),
				),
			},
		},
	})
}

func TestAccOutscaleAccessKey_updatedToActivedKey(t *testing.T) {
	resourceName := "outscale_access_key.outscale_access_key"

	state := "INACTIVE"
	stateUpdated := "ACTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleAccessKeyUpdatedConfig(state),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", state),
				),
			},
			{
				Config: testAccOutscaleAccessKeyUpdatedConfig(stateUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", stateUpdated),
				),
			},
		},
	})
}

func testAccCheckOutscaleAccessKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access ID is set")
		}
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		filter := oscgo.ReadSecretAccessKeyRequest{
			AccessKeyId: rs.Primary.ID,
		}

		_, _, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background()).ReadSecretAccessKeyRequest(filter).Execute()
		if err != nil {
			return fmt.Errorf("Outscale Access Key not found (%s)", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckOutscaleAccessKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_access_key" {
			continue
		}

		filter := oscgo.ReadSecretAccessKeyRequest{
			AccessKeyId: rs.Primary.ID,
		}

		_, _, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background()).ReadSecretAccessKeyRequest(filter).Execute()
		if err != nil {
			return fmt.Errorf("Outscale Access Key still exists (%s)", rs.Primary.ID)
		}
	}
	return nil
}

const testAccOutscaleAccessKeyBasicConfig = `
	resource "outscale_access_key" "outscale_access_key" {
               expiration_date = "2025-01-04"
        }`

func testAccOutscaleAccessKeyUpdatedConfig(state string) string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "outscale_access_key" {
			state = "%s"
                        expiration_date = "2025-08-10T15:19:21"
		}
	`, state)
}
