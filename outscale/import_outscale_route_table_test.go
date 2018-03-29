package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleRouteTable_importBasic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	checkFn := func(s []*terraform.InstanceState) error {
		// Expect 2: group, 1 rules
		if len(s) != 2 {
			return fmt.Errorf("bad states: %#v", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig,
			},

			{
				ResourceName:     "outscale_route_table.foo",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

func TestAccOutscaleRouteTable_complex(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	checkFn := func(s []*terraform.InstanceState) error {
		// Expect 3: group, 2 rules
		if len(s) != 3 {
			return fmt.Errorf("bad states: %#v", s)
		}

		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableConfig_complexImport,
			},

			{
				ResourceName:     "outscale_route_table.mod",
				ImportState:      true,
				ImportStateCheck: checkFn,
			},
		},
	})
}

const testAccRouteTableConfig_complexImport = `
resource "outscale_lin" "default" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tag {
    Name = "tf-rt-import-test"
  }
}

resource "outscale_subnet" "tf_test_subnet" {
  vpc_id                  = "${outscale_lin.default.id}"
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true

  tags {
    Name = "tf-rt-import-test"
  }
}

resource "outscale_public_ip" "nat" {}

resource "outscale_lin_internet_gateway" "gw" {
  vpc_id = "${outscale_lin.default.id}"

  tag {
    Name = "tf-rt-import-test"
  }
}

variable "private_subnet_cidrs" {
  default = "10.0.0.0/24"
}

resource "outscale_nat_gateway" "nat" {
  count         = "${length(split(",", var.private_subnet_cidrs))}"
  allocation_id = "${element(outscale_public_ip.nat.*.id, count.index)}"
  subnet_id     = "${outscale_subnet.tf_test_subnet.id}"
}

resource "outscale_route_table" "mod" {
  count  = "${length(split(",", var.private_subnet_cidrs))}"
  vpc_id = "${outscale_lin.default.id}"

  tag {
    Name = "tf-rt-import-test"
  }

  depends_on = ["outscale_lin_internet_gateway.ogw", "outscale_lin_internet_gateway.gw"]
}

resource "outscale_route" "mod-1" {
  route_table_id         = "${outscale_route_table.mod.id}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${element(outscale_nat_gateway.nat.*.id, count.index)}"
}

resource "outscale_route" "mod" {
  route_table_id            = "${outscale_route_table.mod.id}"
  destination_cidr_block    = "10.181.0.0/16"
}

### vpc bar

resource "outscale_lin" "bar" {
  cidr_block = "10.1.0.0/16"

  tags {
    Name = "tf-rt-import-test"
  }
}

resource "outscale_lin_internet_gateway" "ogw" {
  vpc_id = "${outscale_lin.bar.id}"

  tags {
    Name = "tf-rt-import-test"
  }
}

`
