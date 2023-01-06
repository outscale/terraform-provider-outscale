package outscale

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIFlexibleGpuCatalog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIFlexibleGpuCatalogRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"flexible_gpu_catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"generations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"model_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"v_ram": {
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

func dataSourceOutscaleOAPIFlexibleGpuCatalogRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadFlexibleGpuCatalogRequest{}

	var resp oscgo.ReadFlexibleGpuCatalogResponse
	var err error
	err = resource.Retry(20*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpuCatalog(
			context.Background()).
			ReadFlexibleGpuCatalogRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	fgcs := resp.GetFlexibleGpuCatalog()[:]
	fgc_ret := make([]map[string]interface{}, len(fgcs))

	for k, v := range fgcs {
		n := make(map[string]interface{})
		n["generations"] = v.GetGenerations()
		n["model_name"] = v.GetModelName()
		n["max_cpu"] = v.GetMaxCpu()
		n["max_ram"] = v.GetMaxRam()
		n["v_ram"] = v.GetVRam()
		fgc_ret[k] = n
	}

	if err := d.Set("flexible_gpu_catalog", fgc_ret); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}
