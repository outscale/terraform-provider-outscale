package outscale

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
)

func DataSourceOutscaleRegions() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleRegionsRead,
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

func DataSourceOutscaleRegionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	var resp oscgo.ReadRegionsResponse
	var err error
	var req oscgo.ReadRegionsRequest

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.RegionApi.ReadRegions(context.Background()).ReadRegionsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	regions := resp.GetRegions()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		regs := make([]map[string]interface{}, len(regions))
		for i, region := range regions {
			regs[i] = map[string]interface{}{
				"endpoint":    region.GetEndpoint(),
				"region_name": region.GetRegionName(),
			}
		}

		return set("regions", regs)
	})
}
