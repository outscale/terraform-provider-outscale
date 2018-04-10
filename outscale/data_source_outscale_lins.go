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

func dataSourceOutscaleVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVpcsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpc_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_options_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_tenancy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_set": tagsSchemaComputed(),
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

func dataSourceOutscaleVpcsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcsInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vpcOk := d.GetOk("vpc_id")

	if filtersOk == false && vpcOk == false {
		return fmt.Errorf("filters, or owner must be assigned, or vpc_id(s) must be provided")
	}

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if vpcOk {
		ids := make([]*string, len(v.([]interface{})))

		for k, v := range v.([]interface{}) {
			ids[k] = aws.String(v.(string))
		}

		req.VpcIds = ids
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

	d.SetId(resource.UniqueId())

	vpc_set := make([]map[string]interface{}, len(resp.Vpcs))

	for i, v := range resp.Vpcs {
		vpc := make(map[string]interface{})

		vpc["vpc_id"] = *v.VpcId
		vpc["cidr_block"] = *v.CidrBlock
		vpc["dhcp_options_id"] = *v.DhcpOptionsId
		vpc["instance_tenancy"] = *v.InstanceTenancy
		vpc["state"] = *v.State
		vpc["tag_set"] = tagsToMap(v.Tags)

		vpc_set[i] = vpc
	}

	d.Set("vpc_set", vpc_set)
	d.Set("request_id", resp.RequestId)

	return nil
}
