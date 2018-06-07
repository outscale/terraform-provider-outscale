package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleVpcEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVpcEndpointsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpc_endpoint_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpc_endpoint_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_table_id": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lin_api_access_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_blocks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceOutscaleVpcEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcEndpointsInput{}

	filters, filtersOk := d.GetOk("filter")
	vpcEndpointIDs, vpceIDsOk := d.GetOk("vpc_endpoint_id")

	if filtersOk == false && vpceIDsOk == false {
		return fmt.Errorf("One of filters, or vpc_endpoint_id must be assigned")
	}

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	if vpceIDsOk {
		var ids []*string
		for _, v := range vpcEndpointIDs.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.VpcEndpointIds = ids
	}

	utils.PrintToJSON(req, "VpcEndpoint")

	var err error
	var resp *fcu.DescribeVpcEndpointsOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error

		resp, err = conn.VM.DescribeVpcEndpoints(req)
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
	if resp == nil || len(resp.VpcEndpoints) == 0 {
		return fmt.Errorf("no matching VPC Endpoints found")
	}

	utils.PrintToJSON(resp, "VpcEndpoint Response")

	d.SetId(resource.UniqueId())

	vpcEndpoints := make([]map[string]interface{}, len(resp.VpcEndpoints))

	for k, v := range resp.VpcEndpoints {
		vpce := make(map[string]interface{})

		policy, err := structure.NormalizeJsonString(aws.StringValue(v.PolicyDocument))
		if err != nil {
			return errwrap.Wrapf("policy contains an invalid JSON: {{err}}", err)
		}

		plID, cidrs, err := getPrefixList(conn, aws.StringValue(v.ServiceName))

		if err != nil {
			return err
		}

		vpce["prefix_list_id"] = aws.StringValue(plID)
		vpce["vpc_id"] = aws.StringValue(v.VpcEndpointId)
		vpce["service_name"] = aws.StringValue(v.ServiceName)
		vpce["lin_api_access_id"] = aws.StringValue(v.VpcEndpointId)
		vpce["route_table_id"] = flattenStringList(v.RouteTableIds)
		vpce["policy"] = policy
		vpce["state"] = aws.StringValue(v.State)
		vpce["cidr_blocks"] = cidrs

		vpcEndpoints[k] = vpce
	}

	utils.PrintToJSON(vpcEndpoints, "VPCEndpoints Set")

	d.Set("request_id", resp.RequestId)

	return d.Set("vpc_endpoint_set", vpcEndpoints)
}
