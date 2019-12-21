package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVpnGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lin_to_vpn_gateway_link": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"lin_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleOAPIVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	vpn, vpnOk := d.GetOk("vpn_gateway_id")

	if !filtersOk && !vpnOk {
		return fmt.Errorf("One of vpn_gateway_id or filters must be assigned")
	}

	params := &fcu.DescribeVpnGatewaysInput{}

	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if vpnOk {
		params.VpnGatewayIds = []*string{aws.String(vpn.(string))}
	}

	var resp *fcu.DescribeVpnGatewaysOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnGateways(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}
	if resp == nil || len(resp.VpnGateways) == 0 {
		return fmt.Errorf("no matching VPN gateway found: %#v", params)
	}
	if len(resp.VpnGateways) > 1 {
		return fmt.Errorf("multiple VPN gateways matched; use additional constraints to reduce matches to a single VPN gateway")
	}

	vgw := resp.VpnGateways[0]

	d.SetId(aws.StringValue(vgw.VpnGatewayId))
	vs := make([]map[string]interface{}, len(vgw.VpcAttachments))

	for k, v := range vgw.VpcAttachments {
		vp := make(map[string]interface{})

		vp["state"] = aws.StringValue(v.State)
		vp["lin_id"] = aws.StringValue(v.VpcId)

		vs[k] = vp
	}

	d.Set("lin_to_vpn_gateway_link", vs)
	d.Set("state", aws.StringValue(vgw.State))
	d.Set("type", aws.StringValue(vgw.Type))
	d.Set("tag", tagsToMap(vgw.Tags))
	d.Set("request_id", resp.RequestId)

	return nil
}
