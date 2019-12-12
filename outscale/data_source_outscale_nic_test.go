package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"os"
	"strconv"
	"strings"

	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleOAPIENIDataSource_basic(t *testing.T) {
	var conf oscgo.Nic
	subregion := os.Getenv("OUTSCALE_REGION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_nic.outscale_nic",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPIENIDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleOAPIENIDataSourceConfig(subregion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIENIExists("outscale_nic.outscale_nic", &conf),
					testAccCheckOutscaleOAPIENIAttributes(&conf, subregion),
				),
			},
		},
	})
}

func TestAccOutscaleOAPIENIDataSource_basicFilter(t *testing.T) {
	var conf oscgo.Nic

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

		var resp oscgo.ReadNicsResponse
		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		req := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{NicIds: &[]string{rs.Primary.ID}},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(req)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return err
		}

		if len(resp.GetNics()) != 0 {
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

		conn := testAccProvider.Meta().(*OutscaleClient).OSCAPI
		dnir := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{NicIds: &[]string{rs.Primary.ID}},
		}

		var resp oscgo.ReadNicsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(dnir)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				return nil
			}
			errString := err.Error()
			return fmt.Errorf("Could not find network interface: %s", errString)

		}

		if len(resp.GetNics()) > 0 {
			return fmt.Errorf("Nic with id %s is not destroyed yet", rs.Primary.ID)
		}
	}

	return nil
}

func testAccOutscaleOAPIENIDataSourceConfig(subregion string) string {
	return fmt.Sprintf(`
		resource "outscale_net" "outscale_net" {
			ip_range = "10.0.0.0/16"
			}
			
		resource "outscale_subnet" "outscale_subnet" {
			subregion_name = "%sa"
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
	`, subregion)
}

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
