package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleFlexibleGpuCatalog() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleFlexibleGpuCatalogRead,
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

func DataSourceOutscaleFlexibleGpuCatalogRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadFlexibleGpuCatalogRequest{}

	resp, err := client.ReadFlexibleGpuCatalog(ctx, req, options.WithRetryTimeout(20*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	fgcs := ptr.From(resp.FlexibleGpuCatalog)[:]
	fgc_ret := make([]map[string]any, len(fgcs))

	for k, v := range fgcs {
		n := make(map[string]any)
		n["generations"] = ptr.From(v.Generations)
		n["model_name"] = ptr.From(v.ModelName)
		n["max_cpu"] = ptr.From(v.MaxCpu)
		n["max_ram"] = ptr.From(v.MaxRam)
		n["v_ram"] = ptr.From(v.VRam)
		fgc_ret[k] = n
	}

	if err := d.Set("flexible_gpu_catalog", fgc_ret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}
