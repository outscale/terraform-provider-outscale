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

func dataSourceOutscaleOAPISubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubnetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"sub_region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag": dataSourceTagsSchema(),

			"lin_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"available_ips_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPISubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSubnetsInput{}

	if id := d.Get("subnet_id"); id != "" {
		req.SubnetIds = []*string{aws.String(id.(string))}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *fcu.DescribeSubnetsOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeSubNet(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	if resp == nil || len(resp.Subnets) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	if len(resp.Subnets) > 1 {
		return fmt.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	subnet := resp.Subnets[0]

	d.SetId(*subnet.SubnetId)
	d.Set("subnet_id", subnet.SubnetId)
	d.Set("lin_id", subnet.VpcId)
	d.Set("sub_region_name", subnet.AvailabilityZone)
	d.Set("ip_range", subnet.CidrBlock)
	d.Set("state", subnet.State)
	d.Set("tag", tagsToMap(subnet.Tags))
	d.Set("available_ips_count", subnet.AvailableIpAddressCount)
	d.Set("request_id", resp.RequesterId)

	return nil
}
