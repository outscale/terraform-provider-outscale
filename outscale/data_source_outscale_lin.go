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

func dataSourceOutscaleVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVpcRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
				Optional: true,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag_set": tagsSchemaComputed(),
		},
	}
}

func dataSourceOutscaleVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeVpcsInput{}

	id, ok1 := d.GetOk("vpc_id")
	v, ok2 := d.GetOk("filter")

	if ok1 == false && ok2 == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	if ok1 {
		req.VpcIds = []*string{aws.String(id.(string))}
	}

	if ok2 {
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
	d.Set("vpc_id", vpc.VpcId)
	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("dhcp_options_id", vpc.DhcpOptionsId)
	d.Set("instance_tenancy", vpc.InstanceTenancy)
	d.Set("state", vpc.State)
	d.Set("tag_set", tagsToMap(vpc.Tags))
	d.Set("request_id", resp.RequestId)

	return nil
}
