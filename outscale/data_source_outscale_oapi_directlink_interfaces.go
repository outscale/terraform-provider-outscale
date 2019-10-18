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

func dataSourceOutscaleOAPIDirectLinkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIDirectLinkInterfacesRead,

		Schema: map[string]*schema.Schema{
			"direct_link_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direct_link_interface_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"outscale_private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_asn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bgp_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direct_link_id": {
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
						"direct_link_interface_id": {
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
					},
				},
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

func dataSourceOutscaleOAPIDirectLinkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	var resp *dl.DescribeVirtualInterfacesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
			ConnectionID: aws.String(d.Get("direct_link_id").(string)),
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

	vifs := make([]map[string]interface{}, len(resp.VirtualInterfaces))

	for k, vif := range resp.VirtualInterfaces {
		vifs[k] = map[string]interface{}{
			"bgp_asn":                    aws.Int64Value(vif.Asn),
			"bgp_key":                    aws.StringValue(vif.AuthKey),
			"address_family":             aws.StringValue(vif.AddressFamily),
			"direct_link_id":             aws.StringValue(vif.ConnectionID),
			"client_private_ip":          aws.StringValue(vif.CustomerAddress),
			"site":                       aws.StringValue(vif.Location),
			"account_id":                 aws.StringValue(vif.OwnerAccount),
			"outscale_private_ip":        aws.StringValue(vif.AmazonAddress),
			"vpn_gateway_id":             aws.StringValue(vif.VirtualGatewayID),
			"virtual_interface_id":       aws.StringValue(vif.VirtualInterfaceID),
			"direct_link_Interface_name": aws.StringValue(vif.VirtualInterfaceName),
			"state":                      aws.StringValue(vif.VirtualInterfaceState),
			"type":                       aws.StringValue(vif.VirtualInterfaceType),
			"vlan":                       aws.Int64Value(vif.Vlan),
		}
	}

	d.Set("request_id", resp.RequestID)

	return d.Set("virtual_interfaces", vifs)
}
