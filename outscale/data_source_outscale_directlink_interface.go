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

func dataSourceOutscaleDirectLinkInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleDirectLinkInterfaceRead,

		Schema: map[string]*schema.Schema{
			"virtual_interface_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"amazon_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_interface_type": {
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

func dataSourceOutscaleDirectLinkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	var resp *dl.DescribeVirtualInterfacesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
			VirtualInterfaceID: aws.String(d.Get("virtual_interface_id").(string)),
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

	d.Set("asn", vif.Asn)
	d.Set("auth_key", vif.AuthKey)
	d.Set("address_family", vif.AddressFamily)
	d.Set("connection_id", vif.ConnectionID)
	d.Set("customer_address", vif.CustomerAddress)
	d.Set("location", vif.Location)
	d.Set("owner_account", vif.OwnerAccount)
	d.Set("amazon_address", vif.AmazonAddress)
	d.Set("virtual_gateway_id", vif.VirtualGatewayID)
	d.Set("virtual_interface_id", vif.VirtualInterfaceID)
	d.Set("virtual_interface_name", vif.VirtualInterfaceName)
	d.Set("virtual_interface_state", vif.VirtualInterfaceState)
	d.Set("virtual_interface_type", vif.VirtualInterfaceType)
	d.Set("vlan", vif.Vlan)
	d.Set("request_id", resp.RequestID)

	return nil
}
