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

func DataSourceOutscaleFlexibleGpus() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleFlexibleGpusRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"flexible_gpus": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"model_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flexible_gpu_id": {
							Type:     schema.TypeString,
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

func DataSourceOutscaleFlexibleGpusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	_, IDOk := d.GetOk("flexible_gpu_id")

	if !filtersOk && !IDOk {
		return diag.Errorf("one of filters, or flexible_gpu_id must be assigned")
	}

	var err error
	req := osc.ReadFlexibleGpusRequest{}
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceFlexibleGpuFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadFlexibleGpus(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		errString := err.Error()
		return diag.Errorf("error reading flexible gpu (%s)", errString)
	}

	flexgps := ptr.From(resp.FlexibleGpus)[:]

	if len(flexgps) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	d.SetId(id.UniqueId())

	return diag.FromErr(setOAPIFlexibleGpuAttributes(d, flexgps))
}

func setOAPIFlexibleGpuAttributes(d *schema.ResourceData, fg []osc.FlexibleGpu) error {
	fgpus := make([]map[string]interface{}, len(fg))
	for k, v := range fg {
		fgpu := make(map[string]interface{})

		fgpu["delete_on_vm_deletion"] = v.DeleteOnVmDeletion
		if ptr.From(v.FlexibleGpuId) != "" {
			fgpu["flexible_gpu_id"] = v.FlexibleGpuId
		}
		if ptr.From(v.Generation) != "" {
			fgpu["generation"] = v.Generation
		}
		if ptr.From(v.ModelName) != "" {
			fgpu["model_name"] = v.ModelName
		}
		if ptr.From(v.State) != "" {
			fgpu["state"] = *v.State
		}
		if ptr.From(v.SubregionName) != "" {
			fgpu["subregion_name"] = v.SubregionName
		}
		if ptr.From(v.VmId) != "" {
			fgpu["vm_id"] = v.VmId
		}
		fgpus[k] = fgpu
	}

	return d.Set("flexible_gpus", fgpus)
}
