package oapi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_AccessKey_Basic(t *testing.T) {
	resourceName := "outscale_access_key.basic_access_key"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		PreCheck:                 func() { testacc.PreCheck(t) },
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
	resourceName := "outscale_access_key.update_access_key"
	state := "ACTIVE"
	stateUpdated := "INACTIVE"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		PreCheck:                 func() { testacc.PreCheck(t) },
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		PreCheck:                 func() { testacc.PreCheck(t) },
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		PreCheck:                 func() { testacc.PreCheck(t) },
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

func TestAccOthers_AccessKey_Migration(t *testing.T) {
	state := "INACTIVE"
	stateUpdated := "ACTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps: testacc.FrameworkMigrationTestSteps("1.0.1",
			testAccAccessKeyBasicConfig,
			testAccAccessKeyUpdateState(state),
			testAccAccessKeyUpdateState(stateUpdated),
		),
	})
}

func TestAccOthers_AccessKey_ExpirationDate_Migration(t *testing.T) {
	expirDate := time.Now().AddDate(1, 1, 0).Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps: testacc.FrameworkMigrationTestStepsWithConfigs("1.0.1",
			testacc.MigrationTestConfig{
				Config: testAccAccessKeyExpirationDateConfig(expirDate),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("outscale_access_key.date_access_key", plancheck.ResourceActionUpdate),
					},
				},
			},
		),
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
