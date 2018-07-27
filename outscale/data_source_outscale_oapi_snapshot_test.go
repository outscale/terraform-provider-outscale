package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleOAPISnapshotDataSource_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISnapshotDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotDataSourceID("data.outscale_snapshot.snapshot"),
					resource.TestCheckResourceAttr("data.outscale_snapshot.snapshot", "volume_size", "40"),
				),
			},
		},
	})
}

func TestAccOutscaleOAPISnapshotDataSource_multipleFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOutscaleOAPISnapshotDataSourceConfigWithMultipleFilters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPISnapshotDataSourceID("data.outscale_snapshot.snapshot"),
					resource.TestCheckResourceAttr("data.outscale_snapshot.snapshot", "volume_size", "10"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPISnapshotDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find snapshot data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot data source ID not set")
		}
		return nil
	}
}

const testAccCheckOutscaleOAPISnapshotDataSourceConfig = `
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    type = "gp2"
    size = 40
    tag {
        Name = "External Volume"
    }
}

resource "outscale_snapshot" "snapshot" {
    volume_id = "${outscale_volume.example.id}"
}

data "outscale_snapshot" "snapshot" {
    snapshot_id = "${outscale_snapshot.snapshot.id}"
}
`

const testAccCheckOutscaleOAPISnapshotDataSourceConfigWithMultipleFilters = `
resource "outscale_volume" "external1" {
    availability_zone = "eu-west-2a"
    type = "gp2"
    size = 10
    tag {
        Name = "External Volume 1"
    }
}

resource "outscale_snapshot" "snapshot" {
    volume_id = "${outscale_volume.external1.id}"
}

data "outscale_snapshot" "snapshot" {
    snapshot_id = "${outscale_snapshot.snapshot.id}"
    filter {
	name = "volume-size"
	values = ["10"]
    }
}
`
