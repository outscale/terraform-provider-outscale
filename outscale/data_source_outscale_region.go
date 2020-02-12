package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIRegionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"region_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIRegionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtbID, rtbOk := d.GetOk("region_name")
	filter, filterOk := d.GetOk("filter")

	if !filterOk && !rtbOk {
		return fmt.Errorf("One of region_name or filters must be assigned")
	}

	req := &fcu.DescribeRegionsInput{}

	if rtbOk {
		req.RegionNames = []*string{aws.String(rtbID.(string))}
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
	if len(resp.Regions) > 1 {
		return fmt.Errorf("multiple regions matched; use additional constraints to reduce matches to a single region")
	}

	region := resp.Regions[0]

	d.SetId(*region.RegionName)
	d.Set("region_name", region.RegionName)
	d.Set("region_endpoint", region.Endpoint)
	d.Set("request_id", resp.RequestId)

	return nil
}
