package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleFlexibleGpu() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleFlexibleGpuRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
	}
}

func DataSourceOutscaleFlexibleGpuRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	flexID, IDOk := d.GetOk("flexible_gpu_id")

	if !filtersOk && !IDOk {
		return diag.Errorf("one of filters, or flexible_gpu_id must be assigned")
	}

	var err error
	req := osc.ReadFlexibleGpusRequest{}

	req.Filters = &osc.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{flexID.(string)},
	}

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceFlexibleGpuFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadFlexibleGpus(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.FlexibleGpus == nil || len(*resp.FlexibleGpus) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.FlexibleGpus) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	fg := (*resp.FlexibleGpus)[0]

	if err := d.Set("delete_on_vm_deletion", ptr.From(fg.DeleteOnVmDeletion)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subregion_name", ptr.From(fg.SubregionName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("generation", ptr.From(fg.Generation)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("flexible_gpu_id", ptr.From(fg.FlexibleGpuId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_id", ptr.From(fg.VmId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("model_name", ptr.From(fg.ModelName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", ptr.From(fg.State)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(fg.FlexibleGpuId))
	return nil
}

func buildOutscaleDataSourceFlexibleGpuFilters(set *schema.Set) (*osc.FiltersFlexibleGpu, error) {
	var filters osc.FiltersFlexibleGpu
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "delete_on_vm_deletion":
			filters.DeleteOnVmDeletion = new(cast.ToBool(filterValues[0]))
		case "flexible_gpu_ids":
			filters.FlexibleGpuIds = &filterValues
		case "generations":
			filters.Generations = &filterValues
		case "model_names":
			filters.ModelNames = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.FlexibleGpuState { return osc.FlexibleGpuState(s) }))
		case "subregion_names":
			filters.SubregionNames = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
