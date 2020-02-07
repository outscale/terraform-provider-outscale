package outscale

import (
	"fmt"
	"strings"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIAvailabilityZoneRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"sub_region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIAvailabilityZoneRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	zone, zoneOk := d.GetOk("sub_region_name")

	if !filtersOk && !zoneOk {
		return fmt.Errorf("One of sub_region_name or filters must be assigned")
	}

	req := &fcu.DescribeAvailabilityZonesInput{}

	if zoneOk {
		req.ZoneNames = []*string{aws.String(zone.(string))}
	}

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *fcu.DescribeAvailabilityZonesOutput
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeAvailabilityZones(req)
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
	if resp == nil || len(resp.AvailabilityZones) == 0 {
		return fmt.Errorf("no matching AZ found")
	}
	if len(resp.AvailabilityZones) > 1 {
		return fmt.Errorf("multiple AZs matched; use additional constraints to reduce matches to a single AZ")
	}

	az := resp.AvailabilityZones[0]

	d.SetId(*az.ZoneName)
	d.Set("sub_region_name", az.ZoneName)
	d.Set("region_name", az.RegionName)
	d.Set("state", az.State)
	d.Set("request_id", resp.RequestId)

	return nil
}
