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

func dataSourceOutscaleSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleSubnetsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subnet_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnet_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cidr_block": {
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

						"tag_set": tagsSchemaComputed(),

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"available_ip_address_count": {
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

func dataSourceOutscaleSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSubnetsInput{}

	if id, ok := d.GetOk("subnet_id"); ok {
		var ids []*string
		for _, v := range id.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.SubnetIds = ids
	}

	if filters, filtersOk := d.GetOk("filter"); filtersOk {
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
			subnet["availability_zone"] = *v.AvailabilityZone
		}
		if v.AvailableIpAddressCount != nil {
			subnet["available_ip_address_count"] = *v.AvailableIpAddressCount
		}
		if v.CidrBlock != nil {
			subnet["cidr_block"] = *v.CidrBlock
		}
		if v.State != nil {
			subnet["state"] = *v.State
		}
		if v.SubnetId != nil {
			subnet["subnet_id"] = *v.SubnetId
		}
		if v.Tags != nil {
			subnet["tag_set"] = tagsToMap(v.Tags)
		}
		if v.VpcId != nil {
			subnet["vpc_id"] = *v.VpcId
		}

		subnets[k] = subnet
	}

	if err := d.Set("subnet_set", subnets); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.RequesterId)

	return nil
}
