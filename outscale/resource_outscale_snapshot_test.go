package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleSnapshot_basic(t *testing.T) {
	var v fcu.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("outscale_snapshot.test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleSnapshot_withDescription(t *testing.T) {
	var v fcu.Snapshot
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleSnapshotConfigWithDescription,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSnapshotExists("outscale_snapshot.test", &v),
					resource.TestCheckResourceAttr("outscale_snapshot.test", "description", "Snapshot Acceptance Test"),
				),
			},
		},
	})
}

func testAccCheckSnapshotExists(n string, v *fcu.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		request := &fcu.DescribeSnapshotsInput{
			SnapshotIds: []*string{aws.String(rs.Primary.ID)},
		}

		var resp *fcu.DescribeSnapshotsOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeSnapshots(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err == nil {
			if resp.Snapshots != nil && len(resp.Snapshots) > 0 {
				*v = *resp.Snapshots[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding Snapshot %s", rs.Primary.ID)
	}
}

const testAccOutscaleSnapshotConfig = `
resource "outscale_volume" "test" {
	availability_zone = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
}
`

const testAccOutscaleSnapshotConfigWithDescription = `
resource "outscale_volume" "description_test" {
	availability_zone = "us-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.description_test.id}"
	description = "Snapshot Acceptance Test"
}
`
