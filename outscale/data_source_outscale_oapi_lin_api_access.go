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
)

func dataSourceOutscaleOAPIVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcEndpointRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"lin_api_access_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"prefix_list_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": {
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
			"ip_range": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIVpcEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcEndpointsInput{}

	id, ok1 := d.GetOk("lin_api_access_id")
	v, ok2 := d.GetOk("filter")

	if ok1 == false && ok2 == false {
		return fmt.Errorf("One of filters, or lin_api_access_id must be assigned")
	}

	if ok1 {
		req.VpcEndpointIds = []*string{aws.String(id.(string))}
	}

	if ok2 {
		req.Filters = buildOutscaleDataSourceFilters(v.(*schema.Set))
	}
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
		return fmt.Errorf("no matching VPC found")
	}
	if len(resp.VpcEndpoints) > 1 {
		return fmt.Errorf("multiple VPCs matched; use additional constraints to reduce matches to a single VPC")
	}

	vpc := resp.VpcEndpoints[0]

	policy, err := structure.NormalizeJsonString(aws.StringValue(vpc.PolicyDocument))
	if err != nil {
		return errwrap.Wrapf("policy contains an invalid JSON: {{err}}", err)
	}

	plID, cidrs, err := getOAPIPrefixList(conn, aws.StringValue(vpc.ServiceName))

	if err != nil {
		return err
	}

	if plID != nil {
		d.Set("prefix_list_id", plID)
	}

	d.SetId(*vpc.VpcEndpointId)
	d.Set("vpc_id", vpc.VpcEndpointId)
	d.Set("prefix_list_name", vpc.ServiceName)
	d.Set("route_table_id", flattenStringList(vpc.RouteTableIds))
	d.Set("policy", policy)
	d.Set("state", vpc.State)
	d.Set("ip_ranges", cidrs)
	d.Set("request_id", resp.RequestId)

	return nil
}

func getOAPIPrefixList(conn *fcu.Client, serviceName string) (*string, []interface{}, error) {
	req := &fcu.DescribePrefixListsInput{}
	req.Filters = buildFCUAttributeFilterList(
		map[string]string{
			"prefix-list-name": serviceName,
		},
	)

	var resp *fcu.DescribePrefixListsOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribePrefixLists(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return nil, make([]interface{}, 0), err
	}
	if resp != nil && len(resp.PrefixLists) > 0 {
		if len(resp.PrefixLists) > 1 {
			return nil, make([]interface{}, 0), fmt.Errorf("multiple prefix lists associated with the service name '%s'. Unexpected", serviceName)
		}

		pl := resp.PrefixLists[0]

		return pl.PrefixListId, flattenStringList(pl.Cidrs), nil

	}
	return nil, make([]interface{}, 0), nil

}
