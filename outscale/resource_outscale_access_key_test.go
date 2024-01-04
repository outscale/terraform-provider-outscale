package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOthers_AccessKey_basic(t *testing.T) {
	t.Parallel()
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

func TestAccOthers_AccessKey_updatedToInactivedKey(t *testing.T) {
	t.Parallel()
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

func TestAccOthers_AccessKey_updatedToActivedKey(t *testing.T) {
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

func TestAccOthers_AccessKey_updatedExpirationDate(t *testing.T) {
	resourceName := "outscale_access_key.outscale_access_key"
	expirDate := time.Now().AddDate(1, 1, 0).Format("2006-01-02")
	expirDateUpdated := time.Now().AddDate(1, 4, 0).Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleAccessKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleAccessKeyExpirationDateConfig(string(expirDate)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "expiration_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				Config: testAccOutscaleAccessKeyExpirationDateConfig(string(expirDateUpdated)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleAccessKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "expiration_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
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
		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background()).ReadSecretAccessKeyRequest(filter).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
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
		err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.AccessKeyApi.ReadSecretAccessKey(context.Background()).ReadSecretAccessKeyRequest(filter).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Outscale Access Key not found (%s)", rs.Primary.ID)
		}
	}
	return nil
}

const testAccOutscaleAccessKeyBasicConfig = `
	resource "outscale_access_key" "outscale_access_key" {}`

func testAccOutscaleAccessKeyUpdatedConfig(state string) string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "outscale_access_key" {
			state = "%s"
		}
	`, state)
}

func testAccOutscaleAccessKeyExpirationDateConfig(expirDate string) string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "outscale_access_key" {
			expiration_date = "%s"
		}
	`, expirDate)
}
