package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshot_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapiFlag, err := strconv.ParseBool(o)
	if err != nil {
		oapiFlag = false
	}

	if !oapiFlag {
		t.Skip()
	}
	var v oapi.Snapshots
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshot_withDescription(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapiFlag, err := strconv.ParseBool(o)
	if err != nil {
		oapiFlag = false
	}

	if !oapiFlag {
		t.Skip()
	}
	var v oapi.Snapshots
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISnapshotConfigWithDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOAPISnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func testAccCheckOAPISnapshotExists(n string, v *oapi.Snapshots) resource.TestCheckFunc {
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
			Filters: oapi.Filters_10{SnapshotIds: []string{rs.Primary.ID}},
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

const testAccOutscaleOAPISnapshotConfig = `
resource "outscale_volume" "test" {
	sub_region_name = "us-west-1"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
}
`

const testAccOutscaleOAPISnapshotConfigWithDescription = `
resource "outscale_volume" "description_test" {
	availability_zone = "us-west-1"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.description_test.id}"
	description = "Snapshot Acceptance Test"
}
`
