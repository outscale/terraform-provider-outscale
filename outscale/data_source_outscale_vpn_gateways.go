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

func dataSourceOutscaleOAPIVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpnGatewaysRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_gateway_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_gateway": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_id": {
							Type:     schema.TypeString,
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
						"tag": tagsSchemaComputed(),
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIVpnGatewaysRead(d *schema.ResourceData, meta interface{}) error {
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
		params.VpnGatewayIds = expandStringList(vpn.([]interface{}))
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

	vpns := make([]map[string]interface{}, len(resp.VpnGateways))

	for k, v := range resp.VpnGateways {
		vpn := make(map[string]interface{})

		vs := make([]map[string]interface{}, len(v.VpcAttachments))

		for k, v1 := range v.VpcAttachments {
			vp := make(map[string]interface{})
			vp["state"] = aws.StringValue(v1.State)
			vp["lin_id"] = aws.StringValue(v1.VpcId)

			vs[k] = vp
		}

		vpn["lin_to_vpn_gateway_link"] = vs
		vpn["state"] = aws.StringValue(v.State)
		vpn["vpn_gateway_id"] = aws.StringValue(v.VpnGatewayId)
		vpn["tag"] = tagsToMap(v.Tags)

		vpns[k] = vpn
	}

	d.Set("vpn_gateway", vpns)
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())

	return nil
}
