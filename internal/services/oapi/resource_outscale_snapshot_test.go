package oapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOthers_Snapshot_basic(t *testing.T) {
	var v osc.Snapshot
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(t.Context(), "outscale_snapshot.outscale_snapshot", &v),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_withDescription(t *testing.T) {
	var v osc.Snapshot
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfigWithDescription(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(t.Context(), "outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_CopySnapshot(t *testing.T) {
	var v osc.Snapshot
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfigCopySnapshot(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(t.Context(), "outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Target Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_UpdateTags(t *testing.T) {
	region := utils.GetRegion()
	// var v osc.Snapshot
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfigUpdateTags(region, "Terraform-Snapshot"),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccOutscaleSnapshotConfigUpdateTags(region, "Terraform-Snapshot-2"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccOthers_Snapshot_importBasic(t *testing.T) {
	var v osc.Snapshot
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(t.Context(), "outscale_snapshot.outscale_snapshot", &v),
				),
			},
			testacc.ImportStep("outscale_snapshot.outscale_snapshot", "permissions_to_create_volume", "request_id"),
		},
	})
}

func testAccCheckOAPISnapshotExists(ctx context.Context, n string, v *osc.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no id is set")
		}

		client := testacc.ConfiguredClient.OSC

		request := osc.ReadSnapshotsRequest{
			Filters: &osc.FiltersSnapshot{SnapshotIds: &[]string{rs.Primary.ID}},
		}

		resp, err := client.ReadSnapshots(ctx, request, options.WithRetryTimeout(DefaultTimeout))
		if err == nil {
			if resp.Snapshots != nil && len(*resp.Snapshots) > 0 {
				*v = (*resp.Snapshots)[0]
				return nil
			}
		}
		return fmt.Errorf("error finding snapshot %s", rs.Primary.ID)
	}
}

func testAccOutscaleSnapshotConfig(region string) string {
	return fmt.Sprintf(`
		 resource "outscale_volume" "outscale_volume" {
    subregion_name = "%sa"
    size            = 40
}
resource "outscale_snapshot" "outscale_snapshot" {
    volume_id = outscale_volume.outscale_volume.volume_id
}
resource "outscale_snapshot_attributes" "outscale_snapshot_attributes" {
    snapshot_id = outscale_snapshot.outscale_snapshot.snapshot_id
    permissions_to_create_volume_additions  {
                        account_ids = ["458594607190"]
        }
}
	`, region)
}

func testAccOutscaleSnapshotConfigWithDescription(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "description_test" {
			subregion_name = "%sa"
			size = 1
		}

		resource "outscale_snapshot" "test" {
			volume_id = outscale_volume.description_test.id
			description = "Snapshot Acceptance Test"
		}
	`, region)
}

func testAccOutscaleSnapshotConfigCopySnapshot(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "description_test" {
			subregion_name = "%[1]sb"
			size           = 1
		}

		resource "outscale_snapshot" "source" {
			volume_id   = outscale_volume.description_test.id
			description = "Source Snapshot Acceptance Test"
		}

		resource "outscale_snapshot" "test" {
			source_region_name = "%[1]s"
			source_snapshot_id = outscale_snapshot.source.id
			description        = "Target Snapshot Acceptance Test"
		}
	`, region)
}

func testAccOutscaleSnapshotConfigUpdateTags(region, value string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "outscale_volume" {
		subregion_name = "%sa"
		size           = 10
	  }
	  resource "outscale_snapshot" "outscale_snapshot" {
		volume_id = outscale_volume.outscale_volume.volume_id

		tags {
		  key   = "Name"
		  value = "%s"
		}
	  }
	`, region, value)
}
