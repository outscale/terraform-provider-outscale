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

func dataSourceOutscaleAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleAvailabilityZonesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"zone_name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone_info": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_state": &schema.Schema{
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

func dataSourceOutscaleAvailabilityZonesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	zone, zoneOk := d.GetOk("zone_name")

	if !filtersOk && !zoneOk {
		return fmt.Errorf("One of zone_name or filters must be assigned")
	}

	req := &fcu.DescribeAvailabilityZonesInput{}

	if zoneOk {
		var ids []*string
		for _, v := range zone.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.ZoneNames = ids
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

	d.SetId(resource.UniqueId())

	azi := make([]map[string]interface{}, len(resp.AvailabilityZones))

	for k, v := range resp.AvailabilityZones {
		az := make(map[string]interface{})
		az["region_name"] = *v.RegionName
		az["zone_name"] = *v.ZoneName
		az["zone_state"] = *v.State
		azi[k] = az
	}

	d.Set("availability_zone_info", azi)
	d.Set("request_id", resp.RequestId)

	return nil
}
