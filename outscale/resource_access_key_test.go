package outscale

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOthers_AccessKey_basic(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_access_key.basic_access_key"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyBasicConfig,
				Check: resource.ComposeTestCheckFunc(
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

func TestAccOthers_AccessKeyUpdatedToInactivedKey(t *testing.T) {
	t.Parallel()
	resourceName := "outscale_access_key.update_access_key"
	state := "ACTIVE"
	stateUpdated := "INACTIVE"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyUpdateState(state),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(resourceName, "state", state),
				),
			},
			{
				Config: testAccAccessKeyUpdateState(stateUpdated),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccOthers_AccessKeyUpdatedToActivedKey(t *testing.T) {
	resourceName := "outscale_access_key.update_access_key"

	state := "INACTIVE"
	stateUpdated := "ACTIVE"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyUpdateState(state),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttr(resourceName, "state", state),
				),
			},
			{
				Config: testAccAccessKeyUpdateState(stateUpdated),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccOthers_AccessKeyUpdatedExpirationDate(t *testing.T) {
	resourceName := "outscale_access_key.date_access_key"
	expirDate := time.Now().AddDate(1, 1, 0).Format("2006-01-02")
	expirDateUpdated := time.Now().AddDate(2, 4, 0).Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck:                 func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccAccessKeyExpirationDateConfig(expirDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "access_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(resourceName, "expiration_date"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modification_date"),
					resource.TestCheckResourceAttrSet(resourceName, "secret_key"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				Config: testAccAccessKeyExpirationDateConfig(expirDateUpdated),
				Check: resource.ComposeTestCheckFunc(
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

const testAccAccessKeyBasicConfig = `
	resource "outscale_access_key" "basic_access_key" {}`

func testAccAccessKeyUpdateState(state string) string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "update_access_key" {
			state = "%s"
		}
	`, state)
}

func testAccAccessKeyExpirationDateConfig(expirDate string) string {
	return fmt.Sprintf(`
		resource "outscale_access_key" "date_access_key" {
			expiration_date = "%s"
		}
	`, expirDate)
}
