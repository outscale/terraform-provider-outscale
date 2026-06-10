package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOthers_Snapshot_Basic(t *testing.T) {
	resourceName := "outscale_snapshot.outscale_snapshot"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: snapshotConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
				),
			},
			{
				// Refresh the state after adding a snapshot_attributes to not have an import state drift
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "permissions_to_create_volume.#"),
				),
			},
			testacc.ImportStep(resourceName, testacc.DefaultIgnores()...),
		},
	})
}

func TestAccOthers_Snapshot_CopySnapshot(t *testing.T) {
	resourceName := "outscale_snapshot.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: snapshotConfigCopySnapshot(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "source_snapshot_id"),
				),
			},
			// Ignore optional attributes that cannot be extracted from a snapshot read
			testacc.ImportStep(resourceName, "source_region_name", "source_snapshot_id", "request_id"),
		},
	})
}

func TestAccOthers_Snapshot_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.5.0",
			snapshotConfig(false),
			snapshotConfigCopySnapshot(),
		),
	})
}

func TestAccOthers_Snapshot_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_snapshot.outscale_snapshot"
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-resource-snapshot"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			snapshotConfigWithTag(invalidTagKey, tagValue),
			snapshotConfigWithTag("Name", tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

func snapshotConfig(addPermissions bool) string {
	config := fmt.Sprintf(`
resource "outscale_volume" "outscale_volume" {
    subregion_name = "%sa"
    size            = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
    description = "testacc-snapshot"
}
`, utils.GetRegion())

	if addPermissions {
		config += `
resource "outscale_snapshot_attributes" "outscale_snapshot_attributes" {
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
    permissions_to_create_volume_additions  {
        account_ids = ["458594607190"]
    }
}
`
	}
	return config
}

func snapshotConfigWithTag(tagKey, tagValue string) string {
	return fmt.Sprintf(`
resource "outscale_volume" "outscale_volume" {
    subregion_name = "%sa"
    size           = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
    volume_id   = outscale_volume.outscale_volume.volume_id
    description = "testacc-snapshot"
    tags {
        key   = %q
        value = %q
    }
}
`, utils.GetRegion(), tagKey, tagValue)
}

func snapshotConfigCopySnapshot() string {
	return fmt.Sprintf(`
resource "outscale_volume" "vol" {
	subregion_name = "%[1]sb"
	size           = 1
}

resource "outscale_snapshot" "source" {
	volume_id   = outscale_volume.vol.id
	description = "Source Snapshot Acceptance Test"
}

resource "outscale_snapshot" "test" {
	source_region_name = "%[1]s"
	source_snapshot_id = outscale_snapshot.source.id
	description        = "Target Snapshot Acceptance Test"
}`, utils.GetRegion())
}
