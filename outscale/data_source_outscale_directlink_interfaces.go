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

func dataSourceOutscaleDirectLinkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleDirectLinkInterfacesRead,

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amazon_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_interface_id": {
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

func dataSourceOutscaleDirectLinkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	var err error
	var resp *dl.DescribeVirtualInterfacesOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.API.DescribeVirtualInterfaces(&dl.DescribeVirtualInterfacesInput{
			ConnectionID: aws.String(d.Get("connection_id").(string)),
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
			"asn":                     aws.Int64Value(vif.Asn),
			"auth_key":                aws.StringValue(vif.AuthKey),
			"address_family":          aws.StringValue(vif.AddressFamily),
			"connection_id":           aws.StringValue(vif.ConnectionID),
			"customer_address":        aws.StringValue(vif.CustomerAddress),
			"location":                aws.StringValue(vif.Location),
			"owner_account":           aws.StringValue(vif.OwnerAccount),
			"amazon_address":          aws.StringValue(vif.AmazonAddress),
			"virtual_gateway_id":      aws.StringValue(vif.VirtualGatewayID),
			"virtual_interface_id":    aws.StringValue(vif.VirtualInterfaceID),
			"virtual_interface_name":  aws.StringValue(vif.VirtualInterfaceName),
			"virtual_interface_state": aws.StringValue(vif.VirtualInterfaceState),
			"virtual_interface_type":  aws.StringValue(vif.VirtualInterfaceType),
			"vlan":                    aws.Int64Value(vif.Vlan),
		}
	}

	d.Set("request_id", resp.RequestID)

	return d.Set("virtual_interfaces", vifs)
}
