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

func dataSourceOutscaleOAPISubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPISubnetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subnet_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnet": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tag": tagsSchemaComputed(),

						"lin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"available_ips_count": {
							Type:     schema.TypeInt,
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

func dataSourceOutscaleOAPISubnetsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSubnetsInput{}

	if id := d.Get("subnet_id"); id != "" {
		var ids []*string
		for _, v := range id.([]string) {
			ids = append(ids, aws.String(v))
		}
		req.SubnetIds = ids
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

	subnets := make([]map[string]interface{}, len(resp.Subnets))

	for k, v := range resp.Subnets {
		subnet := make(map[string]interface{})

		if v.AvailabilityZone != nil {
			subnet["sub_region_name"] = *v.AvailabilityZone
		}
		if v.AvailableIpAddressCount != nil {
			subnet["available_ips_count"] = *v.AvailableIpAddressCount
		}
		if v.CidrBlock != nil {
			subnet["ip_range"] = *v.CidrBlock
		}
		if v.State != nil {
			subnet["state"] = *v.State
		}
		if v.SubnetId != nil {
			subnet["subnet_id"] = *v.SubnetId
		}
		if v.Tags != nil {
			subnet["tag"] = tagsToMap(v.Tags)
		}
		if v.VpcId != nil {
			subnet["lin_id"] = *v.VpcId
		}

		subnets[k] = subnet
	}

	if err := d.Set("subnet", subnets); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.RequestId)

	return nil
}
