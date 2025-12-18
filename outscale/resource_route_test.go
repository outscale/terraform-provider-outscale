package outscale

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/outscale/terraform-provider-outscale/utils/testutils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNet_OutscaleRoute_noopdiff(t *testing.T) {
	resourceName := "outscale_route.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteNoopChange,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}

func TestAccNet_ImportRoute_Basic(t *testing.T) {
	resourceName := "outscale_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteNoopChange,
			},
			testutils.ImportStep(resourceName, "request_id", "await_active_state"),
		},
	})
}

func TestAccNet_Route_importWithNatService(t *testing.T) {
	resourceName := "outscale_route.outscale_route_nat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleRouteWithNatService,
			},
			testutils.ImportStep(resourceName, "request_id", "await_active_state", "routes"),
		},
	})
}

func TestAccNet_Route_changeTarget(t *testing.T) {
	resourceName := "outscale_route.rtnatdef"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: computeConfigTestChangeTarget([]string{"nat_service_id"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "nat_service_id"),
				),
			},
			{
				Config: computeConfigTestChangeTarget([]string{"gateway_id"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "gateway_id"),
				),
			},
		},
	})
}

func TestAccNet_Route_onlyOneTarget(t *testing.T) {
	regex := regexp.MustCompile(".*")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:             computeConfigTestChangeTarget([]string{"nat_service_id"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             computeConfigTestChangeTarget([]string{"gateway_id"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             computeConfigTestChangeTarget([]string{"vm_id"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             computeConfigTestChangeTarget([]string{"nic_id"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:             computeConfigTestChangeTarget([]string{"net_peering_id"}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			// net_peering_id with other
			{
				Config:      computeConfigTestChangeTarget([]string{"net_peering_id", "nat_service_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"net_peering_id", "gateway_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"net_peering_id", "vm_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"net_peering_id", "nic_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			// nat_service_id with other
			{
				Config:      computeConfigTestChangeTarget([]string{"nat_service_id", "gateway_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"nat_service_id", "vm_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"nat_service_id", "nic_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			// gateway_id with other
			{
				Config:      computeConfigTestChangeTarget([]string{"gateway_id", "vm_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			// vm_id with other
			{
				Config:      computeConfigTestChangeTarget([]string{"vm_id", "nic_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
			{
				Config:      computeConfigTestChangeTarget([]string{"gateway_id", "nic_id"}),
				PlanOnly:    true,
				ExpectError: regex,
			},
		},
	})
}

func TestAccNet_Route_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: FrameworkMigrationTestSteps("1.1.3",
			testAccOutscaleRouteNoopChange,
			testAccOutscaleRouteWithNatService,
			computeConfigTestChangeTarget([]string{"nat_service_id"}),
		),
	})
}

var testAccOutscaleRouteNoopChange = `
	resource "outscale_net" "test" {
		ip_range = "10.0.0.0/24"
	}

	resource "outscale_route_table" "test" {
		net_id = outscale_net.test.net_id
	}

	resource "outscale_internet_service" "outscale_internet_service" {}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		internet_service_id = outscale_internet_service.outscale_internet_service.id
		net_id              = outscale_net.test.net_id
	}

	resource "outscale_route" "test" {
		gateway_id           = outscale_internet_service.outscale_internet_service.id
		destination_ip_range = "10.0.0.0/16"
		route_table_id       = outscale_route_table.test.route_table_id
	}
`

var testAccOutscaleRouteWithNatService = `
	resource "outscale_net" "outscale_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key   = "name"
			value = "net"
		}
	}

	resource "outscale_subnet" "outscale_subnet" {
		net_id   = outscale_net.outscale_net.net_id
		ip_range = "10.0.0.0/18"
		tags {
			key   = "name"
			value = "subnet"
		}
	}

	resource "outscale_public_ip" "outscale_public_ip" {
		tags {
			key   = "name"
			value = "public_ip"
		}
	}

	resource "outscale_route_table" "outscale_route_table" {
		net_id = outscale_net.outscale_net.net_id
		tags {
			key   = "name"
			value = "route_table"
		}
	}

	resource "outscale_route" "outscale_route" {
		destination_ip_range = "0.0.0.0/0"
		gateway_id           = outscale_internet_service.outscale_internet_service.internet_service_id
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}

	resource "outscale_route_table_link" "outscale_route_table_link" {
		subnet_id      = outscale_subnet.outscale_subnet.subnet_id
		route_table_id = outscale_route_table.outscale_route_table.id
	}

	resource "outscale_internet_service" "outscale_internet_service" {
		tags {
			key   = "name"
			value = "internet_service"
		}
	}

	resource "outscale_internet_service_link" "outscale_internet_service_link" {
		net_id              = outscale_net.outscale_net.net_id
		internet_service_id = outscale_internet_service.outscale_internet_service.id
	}

	resource "outscale_nat_service" "outscale_nat_service" {
		depends_on   = ["outscale_route.outscale_route"]
		subnet_id    = outscale_subnet.outscale_subnet.subnet_id
		public_ip_id = outscale_public_ip.outscale_public_ip.public_ip_id
		tags {
			key   = "name"
			value = "nat"
		}
	}

	resource "outscale_route" "outscale_route_nat" {
		destination_ip_range = "40.0.0.0/16"
		nat_service_id       = outscale_nat_service.outscale_nat_service.nat_service_id
		route_table_id       = outscale_route_table.outscale_route_table.route_table_id
	}
`

func computeConfigTestChangeTarget(targets []string) string {
	var extra_configs []string
	for _, target := range targets {
		switch target {
		case "nat_service_id":
			extra_configs = append(extra_configs, "nat_service_id = outscale_nat_service.nat.nat_service_id")
		case "gateway_id":
			extra_configs = append(extra_configs, "gateway_id = outscale_internet_service.igw.internet_service_id")
		case "vm_id":
			extra_configs = append(extra_configs, "vm_id = \"toto\"")
		case "nic_id":
			extra_configs = append(extra_configs, "nic_id = \"toti\"")
		case "net_peering_id":
			extra_configs = append(extra_configs, "net_peering_id = \"toto\"")
		default:
			extra_configs = append(extra_configs, "")
		}
	}

	return fmt.Sprintf(testAccOutscaleRouteTemplateChangeTarget, strings.Join(extra_configs, "\n"))
}

var testAccOutscaleRouteTemplateChangeTarget = `
resource "outscale_net" "net" {
    ip_range = "10.0.0.0/16"
  tags {
     key = "name"
     value = "netdemo"
    }

}
resource "outscale_internet_service" "igw" {
  tags {
     key = "name"
     value = "igwdemo"
    }
}

resource "outscale_internet_service_link" "igwl" {
    internet_service_id = outscale_internet_service.igw.internet_service_id
    net_id = outscale_net.net.net_id
}

resource "outscale_public_ip" "pub" {
  tags {
     key = "name"
     value = "eipdemo"
    }
}


resource "outscale_route_table" "rtpub" {
  net_id = outscale_net.net.net_id
  tags {
     key = "name"
     value = "rtpub"
    }
}
resource "outscale_route_table" "rtnat" {
  net_id = outscale_net.net.net_id
  tags {
     key = "name"
     value = "rtnat"
    }
}

resource "outscale_subnet" "subnet-pub" {
  net_id   = outscale_net.net.net_id
  ip_range = "10.0.0.0/24"
  tags {
        key   = "Name"
        value = "subnet-pub"
    }
}

# Bind the route table to the public subnet
resource "outscale_route_table_link" "rtblpub" {
	route_table_id = outscale_route_table.rtpub.route_table_id
	subnet_id = outscale_subnet.subnet-pub.subnet_id
}

resource "outscale_route" "rtpubdef" {
	route_table_id = outscale_route_table.rtpub.route_table_id
	 gateway_id = outscale_internet_service.igw.internet_service_id
	destination_ip_range = "0.0.0.0/0"
}

resource "outscale_subnet" "subnet-nat" {
  net_id   = outscale_net.net.net_id
  ip_range = "10.0.1.0/24"
  tags {
        key   = "Name"
        value = "subnet-nat"
    }
}

# Bind the route table to the nat subnet
resource "outscale_route_table_link" "rtblnat" {
	route_table_id = outscale_route_table.rtnat.route_table_id
	subnet_id = outscale_subnet.subnet-nat.subnet_id
}


# Create the NAT gateway once the IGW is bound.
resource "outscale_nat_service" "nat" {
	subnet_id = outscale_subnet.subnet-pub.subnet_id
	public_ip_id = outscale_public_ip.pub.public_ip_id
  depends_on=[outscale_route.rtpubdef, outscale_route_table_link.rtblpub]
}

# Create a NAT route via the NAT gateway
resource "outscale_route" "rtnatdef" {
	route_table_id = outscale_route_table.rtnat.route_table_id
	%v
	destination_ip_range = "0.0.0.0/0"
}
`
