package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"outscale": testAccProvider,
	}

}

func TestGetOMIByRegion(t *testing.T) {
	if omi := getOMIByRegion("eu-west-2", "ubuntu"); omi.OMI != "ami-fbead1f5" {
		t.Fatalf("expected %s, but got %s", "ami-fbead1f5", omi.OMI)
	}
	if omi := getOMIByRegion("eu-west-2", "centos"); omi.OMI != "ami-4a7bf2b3" {
		t.Fatalf("expected %s, but got %s", "ami-4a7bf2b3", omi.OMI)
	}
	if omi := getOMIByRegion("cn-southeast-1", "ubuntu"); omi.OMI != "ami-d0abdc85" {
		t.Fatalf("expected %s, but got %s", "ami-d0abdc85", omi.OMI)
	}
	// default is centos6 eu-west-2
	if omi := getOMIByRegion("", ""); omi.OMI != "ami-4a7bf2b3" {
		t.Fatalf("expected %s, but got %s", "ami-4a7bf2b3", omi.OMI)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func skipIfNoOAPI(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	isOAPI, err := strconv.ParseBool(o)
	if err != nil {
		isOAPI = false
	}

	if !isOAPI {
		t.Skip()
	}
}

func testAccPreCheck(t *testing.T) {
}

func testAccCheckState(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Can't find this `%s` resource", n)
		}

		utils.PrintToJSON(rs, "[Debug] State: \n")

		return nil
	}
}
func testAccWait(n time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(n)
		return nil
	}
}

type Item struct {
	Platform string
	OMI      string
}

func getOMIByRegion(region, platform string) Item {
	if region == "" {
		region = "eu-west-2"
	}
	omis := make(map[string][]Item)
	omis["eu-west-2"] = []Item{Item{Platform: "centos", OMI: "ami-4a7bf2b3"}}
	omis["eu-west-2"] = append(omis["eu-west-2"], Item{Platform: "ubuntu", OMI: "ami-fbead1f5"})

	omis["us-east-2"] = []Item{Item{Platform: "centos", OMI: "ami-8ceca82d"}}
	omis["us-east-2"] = append(omis["us-east-2"], Item{Platform: "ubuntu", OMI: "ami-f2ea59af"})

	omis["us-west-1"] = []Item{Item{Platform: "centos", OMI: "ami-6e94897f"}}
	omis["us-west-1"] = append(omis["us-west-1"], Item{Platform: "ubuntu", OMI: "ami-b1d1f100"})

	omis["cn-southeast-1"] = []Item{Item{Platform: "centos", OMI: "ami-9c559f7b"}}
	omis["cn-southeast-1"] = append(omis["cn-southeast-1"], Item{Platform: "ubuntu", OMI: "ami-d0abdc85"})

	for _, omi := range omis[region] {
		if omi.Platform == platform {
			return omi
		}
	}
	return omis[region][0]
}
