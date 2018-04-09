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

func dataSourceOutscaleOAPIVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"lin_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleOAPIVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcsInput{}

	if id := d.Get("lin_id"); id != "" {
		req.VpcIds = []*string{aws.String(id.(string))}
	}

	if v, ok := d.GetOk("filters"); ok {
		req.Filters = buildOutscaleDataSourceFilters(v.(*schema.Set))
	}

	var err error
	var resp *fcu.DescribeVpcsOutput
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error

		resp, err = conn.VM.DescribeVpcs(req)
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
	if resp == nil || len(resp.Vpcs) == 0 {
		return fmt.Errorf("no matching VPC found")
	}
	if len(resp.Vpcs) > 1 {
		return fmt.Errorf("multiple VPCs matched; use additional constraints to reduce matches to a single VPC")
	}

	vpc := resp.Vpcs[0]

	d.SetId(*vpc.VpcId)
	d.Set("lin_id", vpc.VpcId)
	d.Set("ip_range", vpc.CidrBlock)
	d.Set("dhcp_options_set_id", vpc.DhcpOptionsId)
	d.Set("tenancy", vpc.InstanceTenancy)
	d.Set("state", vpc.State)
	d.Set("tag", tagsToMap(vpc.Tags))
	d.Set("request_id", resp.RequesterId)

	return nil
}
