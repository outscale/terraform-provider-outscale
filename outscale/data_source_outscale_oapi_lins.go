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

func dataSourceOutscaleOAPIVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"lin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_options_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"net_id": {
							Type:     schema.TypeString,
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
						"tag": tagsSchemaComputed(),
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

func dataSourceOutscaleOAPIVpcsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcsInput{}

	filters, filtersOk := d.GetOk("filter")
	v, vpcOk := d.GetOk("net_id")

	if filtersOk == false && vpcOk == false {
		return fmt.Errorf("filters, or owner must be assigned, or net_id(s) must be provided")
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

	lin := make([]map[string]interface{}, len(resp.Vpcs))

	for i, v := range resp.Vpcs {
		vpc := make(map[string]interface{})

		vpc["net_id"] = *v.VpcId
		vpc["ip_range"] = *v.CidrBlock
		vpc["dhcp_options_set_id"] = *v.DhcpOptionsId
		vpc["tenancy"] = *v.InstanceTenancy
		vpc["state"] = *v.State
		vpc["tag"] = tagsToMap(v.Tags)

		lin[i] = vpc
	}

	d.Set("lin", lin)
	d.Set("request_id", resp.RequestId)

	return nil
}
