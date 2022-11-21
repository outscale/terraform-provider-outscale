package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceVMState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceVMStateRead,
		Schema: getVMStateDataSourceSchema(),
	}
}

func getVMStateDataSourceSchema() map[string]*schema.Schema {
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

func dataSourceVMStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if !instanceIDOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := oscgo.ReadVmsStateRequest{}
	if filtersOk {
		params.SetFilters(buildDataSourceVMStateFilters(filters.(*schema.Set)))
	}
	if instanceIDOk {
		filter := oscgo.FiltersVmsState{}
		filter.SetVmIds([]string{instanceID.(string)})
		params.SetFilters(filter)
	}

	params.SetAllVms(false)

	var resp oscgo.ReadVmsStateResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVmsState(context.Background()).ReadVmsStateRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	filteredStates := resp.GetVmStates()[:]

	var state oscgo.VmStates
	if len(filteredStates) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(filteredStates) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	state = filteredStates[0]

	log.Printf("[DEBUG] outscale_vm.state - Single State found: %s", state.GetVmId())
	return vmStateDataAttrSetter(d, &state)
}

func vmStateDataAttrSetter(d *schema.ResourceData, status *oscgo.VmStates) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	d.SetId(status.GetVmId())
	return statusDescriptionVMStateAttributes(setterFunc, status)
}

func statusDescriptionVMStateAttributes(set AttributeSetter, status *oscgo.VmStates) error {

	if err := set("subregion_name", status.GetSubregionName()); err != nil {
		return err
	}
	if err := set("maintenance_events", statusSetVMState(status.GetMaintenanceEvents())); err != nil {
		return err
	}
	if err := set("vm_state", status.GetVmState()); err != nil {
		return err
	}
	if err := set("vm_id", status.GetVmId()); err != nil {
		return err
	}

	return nil
}

func statusSetVMState(status []oscgo.MaintenanceEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(status))
	for k, v := range status {
		s[k] = map[string]interface{}{
			"code":        v.GetCode(),
			"description": v.GetDescription(),
			"not_after":   v.GetNotAfter(),
			"not_before":  v.GetNotBefore(),
		}
	}

	return s
}

func buildDataSourceVMStateFilters(set *schema.Set) oscgo.FiltersVmsState {
	var filters oscgo.FiltersVmsState
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "maintenance_event_codes":
			filters.SetMaintenanceEventCodes(filterValues)
		case "maintenance_event_descriptions":
			filters.SetMaintenanceEventDescriptions(filterValues)
		case "maintenance_events_not_after":
			filters.SetMaintenanceEventsNotAfter(filterValues)
		case "maintenance_events_not_before":
			filters.SetMaintenanceEventsNotBefore(filterValues)
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "vm_ids":
			filters.SetVmIds(filterValues)
		case "vm_states":
			filters.SetVmStates(filterValues)

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
