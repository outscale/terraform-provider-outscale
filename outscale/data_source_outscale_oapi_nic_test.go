package outscale

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func TestAccOutscaleOAPIENIDataSource_basic(t *testing.T) {
	var conf oapi.Nic

	o := os.Getenv("OUTSCALE_OAPI")

	subregion := os.Getenv("OUTSCALE_REGION")
	if subregion == "" {
		subregion = "in-west-2"
	}

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf, subregion),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIENIDataSource_basicFilter(t *testing.T) {
	var conf oapi.Nic

	o := os.Getenv("OUTSCALE_OAPI")

	subregion := os.Getenv("OUTSCALE_REGION")
	if subregion == "" {
		subregion = "in-west-2"
	}

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIDataSourceConfigFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf, subregion),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIENIDataSourceAttributes(conf *fcu.NetworkInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if conf.Attachment != nil {
			return fmt.Errorf("expected attachment to be nil")
		}

		if *conf.AvailabilityZone != "eu-west-2a" {
			return fmt.Errorf("expected availability_zone to be eu-west-2a, but was %s", *conf.AvailabilityZone)
		}

		return nil
	}
}

func testAccCheckOutscaleOAPIENIDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		var resp *oapi.ReadNicsResponse
		var r *oapi.POST_ReadNicsResponses
		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		req := &oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{NicIds: []string{rs.Primary.ID}},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			r, err = conn.POST_ReadNics(*req)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		resp = r.OK

		if err != nil {
			return err
		}

		if len(resp.Nics) != 0 {
			return fmt.Errorf("Nic is not destroyed yet")
		}
	}
	return nil
}

func testAccCheckOutscaleOAPINICDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_nic" {
			continue
		}

		conn := testAccProvider.Meta().(*OutscaleClient).OAPI
		dnir := &oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{NicIds: []string{rs.Primary.ID}},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(*dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || describeResp.OK == nil {
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
					return nil
				}
				errString = err.Error()
			} else if describeResp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
			} else if describeResp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
			} else if describeResp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
			}
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		if len(describeResp.OK.Nics) > 0 {
			return fmt.Errorf("Nic with id %s is not destroyed yet", rs.Primary.ID)
		}
	}

	return nil
}

const testAccOutscaleOAPIENIDataSourceConfig = `
resource "outscale_net" "outscale_net" {
	ip_range = "10.0.0.0/16"
  }
  
resource "outscale_subnet" "outscale_subnet" {
	subregion_name = "eu-west-2a"
	ip_range       = "10.0.0.0/16"
	net_id         = "${outscale_net.outscale_net.id}"
}

resource "outscale_nic" "outscale_nic" {
	subnet_id = "${outscale_subnet.outscale_subnet.id}"
	tags {
		value = "tf-value"
		key   = "tf-key"
	}
}

data "outscale_nic" "outscale_nic" {
	nic_id = "${outscale_nic.outscale_nic.id}"
}  
`

const testAccOutscaleOAPIENIDataSourceConfigFilter = `
resource "outscale_net" "outscale_net" {
	ip_range = "10.0.0.0/16"
  }
  
resource "outscale_subnet" "outscale_subnet" {
	subregion_name = "eu-west-2a"
	ip_range       = "10.0.0.0/16"
	net_id         = "${outscale_net.outscale_net.id}"
}

resource "outscale_nic" "outscale_nic" {
	subnet_id = "${outscale_subnet.outscale_subnet.id}"
	tags {
		value = "tf-value"
		key   = "tf-key"
	}
}

data "outscale_nic" "outscale_nic" {
	filter {
        name = "nic_ids"
        values = ["${outscale_nic.outscale_nic.nic_id}"]
    } 
}  
`
