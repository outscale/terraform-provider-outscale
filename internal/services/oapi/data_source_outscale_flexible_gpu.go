package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleFlexibleGpu() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleFlexibleGpuRead,
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

func DataSourceOutscaleFlexibleGpuRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	flexID, IDOk := d.GetOk("flexible_gpu_id")

	if !filtersOk && !IDOk {
		return fmt.Errorf("One of filters, or flexible_gpu_id must be assigned")
	}

	var err error
	req := oscgo.ReadFlexibleGpusRequest{}

	req.Filters = &oscgo.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{flexID.(string)},
	}

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourceFlexibleGpuFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadFlexibleGpusResponse

	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(
			context.Background()).ReadFlexibleGpusRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetFlexibleGpus()), "FlexibleGpu"); err != nil {
		return err
	}

	fg := (*resp.FlexibleGpus)[0]

	if err := d.Set("delete_on_vm_deletion", fg.GetDeleteOnVmDeletion()); err != nil {
		return err
	}
	if err := d.Set("subregion_name", fg.GetSubregionName()); err != nil {
		return err
	}
	if err := d.Set("generation", fg.GetGeneration()); err != nil {
		return err
	}
	if err := d.Set("flexible_gpu_id", fg.GetFlexibleGpuId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", fg.GetVmId()); err != nil {
		return err
	}
	if err := d.Set("model_name", fg.GetModelName()); err != nil {
		return err
	}
	if err := d.Set("state", fg.GetState()); err != nil {
		return err
	}
	d.SetId(fg.GetFlexibleGpuId())
	return nil
}

func buildOutscaleDataSourceFlexibleGpuFilters(set *schema.Set) (*oscgo.FiltersFlexibleGpu, error) {
	var filters oscgo.FiltersFlexibleGpu
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "delete_on_vm_deletion":
			filters.SetDeleteOnVmDeletion(cast.ToBool(filterValues[0]))
		case "flexible_gpu_ids":
			filters.SetFlexibleGpuIds(filterValues)
		case "generations":
			filters.SetGenerations(filterValues)
		case "model_names":
			filters.SetModelNames(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "vm_ids":
			filters.SetVmIds(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
