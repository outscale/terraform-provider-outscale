package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
)

func dataSourceOutscaleOAPIDirectLinkInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIDirectLinkInterfaceRead,

		Schema: map[string]*schema.Schema{
			"direct_link_interface_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"outscale_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direct_link_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bgp_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"site": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direct_link_interface_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func dataSourceOutscaleOAPIDirectLinkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	var resp *dl.DescribeVirtualInterfacesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
			VirtualInterfaceID: aws.String(d.Get("direct_link_interface_id").(string)),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading Direct Connect virtual interface: %s", err)
	}

	if resp == nil {
		log.Printf("[WARN] Direct Connect virtual interface (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	vif := resp.VirtualInterfaces[0]

	d.Set("bgp_asn", vif.Asn)
	d.Set("bgp_key", vif.AuthKey)
	d.Set("address_family", vif.AddressFamily)
	d.Set("direct_link_id", vif.ConnectionID)
	d.Set("client_private_ip", vif.CustomerAddress)
	d.Set("site", vif.Location)
	d.Set("account_id", vif.OwnerAccount)
	d.Set("outscale_private_ip", vif.AmazonAddress)
	d.Set("vpn_gateway_id", vif.VirtualGatewayID)
	d.Set("direct_link_interface_id", vif.VirtualInterfaceID)
	d.Set("direct_link_interface_name", vif.VirtualInterfaceName)
	d.Set("state", vif.VirtualInterfaceState)
	d.Set("type", vif.VirtualInterfaceType)
	d.Set("vlan", vif.Vlan)
	d.Set("request_id", resp.RequestID)

	return nil
}
