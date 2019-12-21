package outscale

import (
	"context"
	"errors"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceOutscaleOAPIVMState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMStateRead,
		Schema: getOAPIVMStateDataSourceSchema(),
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

func dataSourceOutscaleOAPIVMStateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if !instanceIDOk && !filtersOk {
		return errors.New("vm_id or filter must be set")
	}

	params := oscgo.ReadVmsStateRequest{}
	if filtersOk {
		params.SetFilters(buildOutscaleOAPIDataSourceVMStateFilters(filters.(*schema.Set)))
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
		resp, _, err = conn.VmApi.ReadVmsState(context.Background(), &oscgo.ReadVmsStateOpts{ReadVmsStateRequest: optional.NewInterface(params)})
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

	log.Printf("[DEBUG] outscale_oapi_vm_state - Single State found: %s", state.GetVmId())
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	return vmStateDataAttrSetter(d, &state)
}

func vmStateDataAttrSetter(d *schema.ResourceData, status *oscgo.VmStates) error {
	setterFunc := func(key string, value interface{}) error {
		return d.Set(key, value)
	}
	d.SetId(status.GetVmId())
	return statusDescriptionOAPIVMStateAttributes(setterFunc, status)
}

func statusDescriptionOAPIVMStateAttributes(set AttributeSetter, status *oscgo.VmStates) error {

	set("subregion_name", status.GetSubregionName())
	set("maintenance_events", statusSetOAPIVMState(status.GetMaintenanceEvents()))
	set("vm_state", status.GetVmState())
	set("vm_id", status.GetVmId())

	return nil
}

func statusSetOAPIVMState(status []oscgo.MaintenanceEvent) []map[string]interface{} {
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

func buildOutscaleOAPIDataSourceVMStateFilters(set *schema.Set) oscgo.FiltersVmsState {
	var filters oscgo.FiltersVmsState
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "maintenance_event_codes":
			//filters.MaintenanceEventCodes = filterValues
		case "maintenance_event_descriptions":
			//filters.MaintenanceEventDescriptions = filterValues
		case "maintenance_events_not_after":
			//filters.MaintenanceEventsNotAfter = filterValues
		case "maintenance_events_not_before":
			//filters.MaintenanceEventsNotBefore = filterValues
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
