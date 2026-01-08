package oapi

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleFlexibleGpuCatalog() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleFlexibleGpuCatalogRead,
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

func DataSourceOutscaleFlexibleGpuCatalogRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	req := oscgo.ReadFlexibleGpuCatalogRequest{}

	var resp oscgo.ReadFlexibleGpuCatalogResponse
	var err = retry.Retry(20*time.Second, func() *retry.RetryError {
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

	d.SetId(id.UniqueId())

	return nil
}
