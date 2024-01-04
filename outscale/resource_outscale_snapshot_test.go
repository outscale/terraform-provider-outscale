package outscale

import (
	"context"
	"fmt"
	"testing"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOthers_Snapshot_basic(t *testing.T) {
	t.Parallel()

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.outscale_snapshot", &v),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_withDescription(t *testing.T) {
	t.Parallel()

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigWithDescription(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_CopySnapshot(t *testing.T) {
	t.Parallel()

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigCopySnapshot(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Target Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func TestAccOthers_Snapshot_UpdateTags(t *testing.T) {
	t.Parallel()

	region := utils.GetRegion()
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

func TestAccOthers_Snapshot_importBasic(t *testing.T) {
	t.Parallel()

	var v oscgo.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig(utils.GetRegion()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.outscale_snapshot", &v),
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
			volume_id = outscale_volume.description_test.id
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

func testAccOutscaleOAPISnapshotConfigUpdateTags(region, value string) string {
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
