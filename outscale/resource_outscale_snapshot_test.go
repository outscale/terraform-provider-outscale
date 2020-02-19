package outscale

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOutscaleOAPISnapshot_basic(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.outscale_snapshot", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshot_withDescription(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigCopySnapshot(region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Target Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshot_UpdateTags(t *testing.T) {
	region := os.Getenv("OUTSCALE_REGION")

	//var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigUpdateTags(region, "Terraform-Snapshot"),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccOutscaleOAPISnapshotConfigUpdateTags(region, "Terraform-Snapshot-2"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
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
			resp, _, err = conn.SnapshotApi.ReadSnapshots(context.Background(), &oscgo.ReadSnapshotsOpts{ReadSnapshotsRequest: optional.NewInterface(request)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
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

func testAccOutscaleOAPISnapshotConfig(region string) string {
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
			source_region_name = "%[1]s"
			source_snapshot_id = "${outscale_snapshot.source.id}"
			description        = "Target Snapshot Acceptance Test"
		}
	`, region)
}

func testAccOutscaleOAPISnapshotConfigUpdateTags(region, value string) string {
	return fmt.Sprintf(`
	resource "outscale_volume" "outscale_volume" {
		subregion_name = "%sa"
		size           = 10
	  }
	  resource "outscale_snapshot" "outscale_snapshot" {
		volume_id = "${outscale_volume.outscale_volume.volume_id}"
		
		tags {
		  key   = "Name"
		  value = "%s"
		}
	  }	  
	`, region, value)
}
