package outscale

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func dataSourceOutscaleOAPIRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIRegionsRead,
		Schema: map[string]*schema.Schema{
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOutscaleOAPIRegionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadRegionsResponse
	var err error
	var req oscgo.ReadRegionsRequest

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.RegionApi.ReadRegions(context.Background()).ReadRegionsRequest(req).Execute()
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

	regions := resp.GetRegions()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		regs := make([]map[string]interface{}, len(regions))
		for i, region := range regions {
			regs[i] = map[string]interface{}{
				"endpoint":    region.GetEndpoint(),
				"region_name": region.GetRegionName(),
			}
		}

		if err := d.Set("request_id", resp.ResponseContext.RequestId); err != nil {
			return err
		}

		return set("regions", regs)
	})
}
