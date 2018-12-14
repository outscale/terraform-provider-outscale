package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIVMState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMStateRead,
		Schema: getOAPIVMStateDataSourceSchema(),
	}
}

func dataSourceOutscaleOAPIVMStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	instanceIds, instanceIdsOk := d.GetOk("vm_id")

	if !instanceIdsOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := oapi.ReadVmsStateRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVmStateFilters(filters.(*schema.Set))
	}
	if instanceIdsOk {
		var ids []string

		for _, id := range instanceIds.(*schema.Set).List() {
			ids = append(ids, id.(string))
		}

		params.Filters.VmIds = ids
	}

	params.AllVms = false

	var resp *oapi.POST_ReadVmsStateResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadVmsState(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	filteredStates := resp.OK.VmStates[:]

	var state oapi.VmStates
	if len(filteredStates) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(filteredStates) > 1 {
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	state = filteredStates[0]

	log.Printf("[DEBUG] outscale_oapi_vm_state - Single State found: %s", state.VmId)

	return vmStateDataAttrSetter(d, &state)
}

func vmStateDataAttrSetter(d *schema.ResourceData, status *oapi.VmStates) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	d.SetId(status.VmId)
	return statusDescriptionOAPIVMStateAttributes(setterFunc, status)
}

func statusDescriptionOAPIVMStateAttributes(set AttributeSetter, status *oapi.VmStates) error {

	set("subregion_name", status.SubregionName)
	set("maintenance_events", oapiEventsSet(status.MaintenanceEvents))
	set("vm_state", status.VmState)
	set("vm_id", status.VmId)
	set("request_id", status.VmState)

	return nil
}

func statusSetOAPIVMState(status []oapi.MaintenanceEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(status))
	for k, v := range status {
		s[k] = map[string]interface{}{
			"code":        v.Code,
			"description": v.Description,
			"not_after":   v.NotAfter,
			"not_before":  v.NotBefore,
		}
	}

	return s
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

func buildOutscaleOAPIDataSourceVmStateFilters(set *schema.Set) oapi.FiltersVmsState {
	var filters oapi.FiltersVmsState
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "maintenance_event_codes":
			filters.MaintenanceEventCodes = filterValues
		case "maintenance_event_descriptions":
			filters.MaintenanceEventDescriptions = filterValues
		case "maintenance_events_not_after":
			filters.MaintenanceEventsNotAfter = filterValues
		case "maintenance_events_not_before":
			filters.MaintenanceEventsNotBefore = filterValues
		case "subregion_names":
			filters.SubregionNames = filterValues
		case "vm_ids":
			filters.VmIds = filterValues
		case "vm_states":
			filters.VmStates = filterValues

		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}

func oapiEventsSet(events []oapi.MaintenanceEvent) []map[string]interface{} {
	s := make([]map[string]interface{}, len(events))

	for k, v := range events {
		status := map[string]interface{}{
			"code":        v.Code,
			"description": v.Description,
			"not_before":  v.NotBefore,
			"not_after":   v.NotAfter,
		}
		s[k] = status
	}
	return s
}
