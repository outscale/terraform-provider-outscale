package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIDLVPNGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIDLVPNGatewaysRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_ids": {
				Type:     schema.TypeString,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": dataSourceFiltersSchema(),
			"virtual_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_gateway_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_gateway_state": &schema.Schema{
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
	}
}

func dataSourceOutscaleOAPIDLVPNGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).DL

	request := &dl.DescribeVirtualGatewaysInput{}

	var getResp *dl.DescribeVirtualGatewaysOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		getResp, err = conn.API.DescribeVirtualGateways(request)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "NoSuchEntity") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading sites: %s", err)
	}

	gateways := make([]map[string]interface{}, len(getResp.VirtualGateways))

	for k, v := range getResp.VirtualGateways {
		gateway := make(map[string]interface{})
		gateway["virtual_gateway_id"] = aws.StringValue(v.VirtualGatewayID)
		gateway["virtual_gateway_state"] = aws.StringValue(v.VirtualGatewayState)

		gateways[k] = gateway
	}

	d.SetId(resource.UniqueId())
	d.Set("virtual_gateways", gateways)

	return d.Set("request_id", getResp.RequestID)
}
