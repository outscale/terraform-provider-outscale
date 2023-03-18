package outscale

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceSnapshot_complete(t *testing.T) {
	t.Parallel()
	region := os.Getenv("OUTSCALE_REGION")
	resourceName := "outscale_snapshot.outscale_snapshot"

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig(region, "Terraform-Snapshot"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", "Snapshot Acceptance Test"),
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume_global_permission", "false"),
				),
			},
			{
				Config: testAccOutscaleOAPISnapshotUpdateWithCopy(region, "Terraform-Snapshot-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.outscale_snapshot_copy", &v),
					testAccCheckOAPISnapshotExists(resourceName, &v),
					resource.TestCheckResourceAttr("outscale_snapshot.outscale_snapshot_copy", "description", "Target Snapshot Acceptance Test"),
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume_global_permission", "true"),
				),
			},
			{
				Config: testAccOutscaleOAPISnapshotUpdate(region, "Terraform-Snapshot-2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "permissions_to_create_volume_global_permission", "false"),
				),
			},
			{
				ResourceName:            "outscale_snapshot.outscale_snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func testAccOutscaleOAPISnapshotConfig(region, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%sa"
			size            = 1
		}

		resource "outscale_snapshot" "outscale_snapshot" {
    		volume_id = outscale_volume.outscale_volume.volume_id
			description = "Snapshot Acceptance Test"
			tags {
				key   = "Name"
				value = "%s"
			}
		}
	`, region, tag)
}

func testAccOutscaleOAPISnapshotUpdateWithCopy(region, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%sa"
			size            = 1
		}

		resource "outscale_snapshot" "outscale_snapshot" {
    		volume_id = outscale_volume.outscale_volume.volume_id
			description = "Snapshot Acceptance Test"
			permissions_to_create_volume_global_permission = true
			permissions_to_create_volume_account_ids = ["458594607190", "458594607191"]
			tags {
				key   = "Name"
				value = "%s"
			}
		}

		resource "outscale_snapshot" "outscale_snapshot_copy" {
			source_region_name = "%[1]s"
			source_snapshot_id = "${outscale_snapshot.outscale_snapshot.id}"
			description        = "Target Snapshot Acceptance Test"
		}

	`, region, tag)
}

func testAccOutscaleOAPISnapshotUpdate(region, tag string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "outscale_volume" {
			subregion_name = "%sa"
			size            = 1
		}

		resource "outscale_snapshot" "outscale_snapshot" {
    		volume_id = outscale_volume.outscale_volume.volume_id
			description = "Snapshot Acceptance Test"
			permissions_to_create_volume_global_permission = false
			permissions_to_create_volume_account_ids = ["458594607192", "458594607191"]
			tags {
				key   = "Name"
				value = "%s"
			}
		}
	`, region, tag)
}

func testAccCheckOAPISnapshotExists(n string, v *oscgo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI

		request := oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadSnapshotsResponse
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			rp, httpResp, err := conn.SnapshotApi.ReadSnapshots(context.Background()).ReadSnapshotsRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err == nil {
			if resp.GetSnapshots() != nil && len(resp.GetSnapshots()) > 0 {
				*v = resp.GetSnapshots()[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding Snapshot %s", rs.Primary.ID)
	}
}
