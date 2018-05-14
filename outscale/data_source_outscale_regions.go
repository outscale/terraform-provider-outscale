package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleRegionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_info": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_endpoint": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleRegionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtbID, rtbOk := d.GetOk("region_name")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !rtbOk {
		return fmt.Errorf("One of region_name or filters must be assigned")
	}

	req := &fcu.DescribeRegionsInput{}

	if rtbOk {
		var ids []*string
		for _, v := range rtbID.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.RegionNames = ids
	}

	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

	log.Printf("[DEBUG] DescribeRegions %+v\n", req)

	var resp *fcu.DescribeRegionsOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeRegions(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if resp == nil || len(resp.Regions) == 0 {
		return fmt.Errorf("no matching regions found")
	}

	d.SetId(resource.UniqueId())

	ri := make([]map[string]interface{}, len(resp.Regions))

	for k, v := range resp.Regions {
		r := make(map[string]interface{})
		r["region_endpoint"] = *v.Endpoint
		r["region_name"] = *v.RegionName
		ri[k] = r
	}

	d.Set("region_info", ri)
	d.Set("request_id", resp.RequestId)

	return nil
}
