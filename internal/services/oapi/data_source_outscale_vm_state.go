package oapi

import (
	"context"
	"log"
	"time"

	"github.com/oapi-codegen/runtime/types"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/iso8601"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/framework/fwhelpers/from"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVMState() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVMStateRead,
		Schema:      getOAPIVMStateDataSourceSchema(),
	}
}

func getOAPIVMStateDataSourceSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
	}

	for k, v := range getVMStateAttrsSchema() {
		wholeSchema[k] = v
	}

	wholeSchema["request_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return wholeSchema
}

func getVMStateAttrsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subregion_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"maintenance_events": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_after": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"not_before": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"all_vms": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vm_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceOutscaleVMStateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if !instanceIDOk && !filtersOk {
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

	if instanceIDOk {
		filter := osc.FiltersVmsState{}
		filter.VmIds = &[]string{instanceID.(string)}
		params.Filters = &filter
	}
	params.AllVms = new(d.Get("all_vms").(bool))

	resp, err := client.ReadVmsState(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	filteredStates := ptr.From(resp.VmStates)[:]

	var state osc.VmStates
	if len(filteredStates) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	if len(filteredStates) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	state = filteredStates[0]

	log.Printf("[DEBUG] outscale_oapi_vm_state - Single State found: %s", state.VmId)
	return diag.FromErr(vmStateDataAttrSetter(d, &state))
}

func vmStateDataAttrSetter(d *schema.ResourceData, status *osc.VmStates) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	d.SetId(status.VmId)
	return statusDescriptionOAPIVMStateAttributes(setterFunc, status)
}

func statusDescriptionOAPIVMStateAttributes(set AttributeSetter, status *osc.VmStates) error {
	if err := set("subregion_name", status.SubregionName); err != nil {
		return err
	}
	if err := set("maintenance_events", statusSetOAPIVMState(status.MaintenanceEvents)); err != nil {
		return err
	}
	if err := set("vm_state", status.VmState); err != nil {
		return err
	}
	if err := set("vm_id", status.VmId); err != nil {
		return err
	}

	return nil
}

func statusSetOAPIVMState(status []osc.MaintenanceEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(status))
	for k, v := range status {
		s[k] = map[string]interface{}{
			"code":        v.Code,
			"description": v.Description,
			"not_after":   from.ISO8601(v.NotAfter),
			"not_before":  from.ISO8601(v.NotBefore),
		}
	}

	return s
}

func buildOutscaleDataSourceVMStateFilters(set *schema.Set) (*osc.FiltersVmsState, error) {
	var filters osc.FiltersVmsState
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "maintenance_event_codes":
			filters.MaintenanceEventCodes = &filterValues
		case "maintenance_event_descriptions":
			filters.MaintenanceEventDescriptions = &filterValues
		case "maintenance_events_not_after":
			var events []types.Date
			for _, s := range filterValues {
				t, err := iso8601.ParseString(s)
				if err != nil {
					return nil, err
				}
				events = append(events, types.Date{Time: t.Time})
			}
			filters.MaintenanceEventsNotAfter = &events
		case "maintenance_events_not_before":
			var events []types.Date
			for _, s := range filterValues {
				t, err := iso8601.ParseString(s)
				if err != nil {
					return nil, err
				}
				events = append(events, types.Date{Time: t.Time})
			}
			filters.MaintenanceEventsNotBefore = &events
		case "subregion_names":
			filters.SubregionNames = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
		case "vm_states":
			filters.VmStates = new(lo.Map(filterValues, func(s string, _ int) osc.VmState {
				return osc.VmState(s)
			}))

		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
