package outscale

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshot_basic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oapi.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshot_withDescription(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oapi.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigWithDescription(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshot_CopySnapshot(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oapi.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigCopySnapshot(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func testAccCheckOAPISnapshotExists(n string, v *oapi.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI

		request := oapi.ReadSnapshotsRequest{
			Filters: oapi.FiltersSnapshot{SnapshotIds: []string{rs.Primary.ID}},
		}

		var resp *oapi.POST_ReadSnapshotsResponses
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.POST_ReadSnapshots(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err == nil {
			if resp.OK.Snapshots != nil && len(resp.OK.Snapshots) > 0 {
				*v = resp.OK.Snapshots[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding Snapshot %s", rs.Primary.ID)
	}
}

func testAccOutscaleOAPISnapshotConfig(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "test" {
			subregion_name = "%sa"
			size = 1
		}

		resource "outscale_snapshot" "test" {
			volume_id = "${outscale_volume.test.id}"
		}
	`, region)
}

func testAccOutscaleOAPISnapshotConfigWithDescription(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "description_test" {
			subregion_name = "%sa"
			size = 1
		}

		resource "outscale_snapshot" "test" {
			volume_id = "${outscale_volume.description_test.id}"
			description = "Snapshot Acceptance Test"
		}
	`, region)
}

func testAccOutscaleOAPISnapshotConfigCopySnapshot(region string) string {
	return fmt.Sprintf(`
		resource "outscale_volume" "description_test" {
			subregion_name = "%[1]sb"
			size           = 1
		}

		resource "outscale_snapshot" "source" {
			volume_id   = "${outscale_volume.description_test.id}"
			description = "Source Snapshot Acceptance Test"
		}

		resource "outscale_snapshot" "test" {
			source_region_name = "%[1]sa"
			source_snapshot_id = "${outscale_snapshot.source.id}"
			description        = "Target Snapshot Acceptance Test"
		}
	`, region)
}
