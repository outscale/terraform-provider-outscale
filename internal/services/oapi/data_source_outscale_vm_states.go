package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVMStates() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVMStatesRead,
		Schema:      getOAPIVMStatesDataSourceSchema(),
	}
}

func getOAPIVMStatesDataSourceSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"all_vms": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"vm_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vm_states": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getVMStateAttrsSchema(),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return wholeSchema
}

func DataSourceOutscaleVMStatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_ids")

	if !instanceIdsOk && !filtersOk {
		return diag.Errorf("vm_id or filter must be set")
	}

	var err error
	params := osc.ReadVmsStateRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVMStateFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if instanceIdsOk {
		filter := osc.FiltersVmsState{}
		filter.VmIds = new(utils.InterfaceSliceToStringSlice(instanceIds.([]interface{})))
		params.Filters = &filter
	}
	params.AllVms = new(d.Get("all_vms").(bool))
	resp, err := client.ReadVmsState(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	filteredStates := ptr.From(resp.VmStates)[:]

	if len(filteredStates) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	return diag.FromErr(statusDescriptionOAPIVMStatesAttributes(d, filteredStates))
}

func statusDescriptionOAPIVMStatesAttributes(d *schema.ResourceData, status []osc.VmStates) error {
	d.SetId(id.UniqueId())

	states := make([]map[string]interface{}, len(status))

	for k, v := range status {
		state := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			state[key] = value
			return nil
		}

		if err := statusDescriptionOAPIVMStateAttributes(setterFunc, &v); err != nil {
			return err
		}

		states[k] = state
	}

	return d.Set("vm_states", states)
}
