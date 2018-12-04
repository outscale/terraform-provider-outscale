package outscale

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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

	return statusDescriptionOAPIVMStateAttributes(d, state)
}

func statusDescriptionOAPIVMStateAttributes(d *schema.ResourceData, status oapi.VmStates) error {

	d.SetId(status.VmId)

	d.Set("sub_region_name", status.SubregionName)

	events := oapiEventsSet(status.MaintenanceEvents)
	err := d.Set("maintenance_event", events)
	if err != nil {
		return err
	}

	err = d.Set("state", status.VmState)
	if err != nil {
		return err
	}

	err = d.Set("comment_item", status.VmState)
	if err != nil {
		return err
	}

	d.Set("comment_state", status.VmState)

	return nil
}

func statusSetOAPIVMState(status *fcu.InstanceStatusSummary) map[string]interface{} {

	st := map[string]interface{}{
		"state": *status.Status,
		"item":  detailsSetOAPIVMState(status.Details),
	}

	return st
}

func detailsSetOAPIVMState(details []*fcu.InstanceStatusDetails) []map[string]interface{} {
	s := make([]map[string]interface{}, len(details))

	for k, v := range details {

		status := map[string]interface{}{
			"name":  *v.Name,
			"state": *v.Status,
		}
		s[k] = status
	}

	return s
}

func flattenedStateOAPIVMState(state *fcu.InstanceState) map[string]interface{} {
	return map[string]interface{}{
		"code": fmt.Sprintf("%d", *state.Code),
		"name": *state.Name,
	}
}

func getOAPIVMStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"sub_region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"maintenance_event": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state_code": {
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
		"state": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state_code": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"comment_item": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"comment_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
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
		case "maintenance­-event-code":
			filters.MaintenanceEventCodes = filterValues
		case "maintenance­-event-description":
			filters.MaintenanceEventDescriptions = filterValues
		case "maintenance­-event-not­after":
			filters.MaintenanceEventsNotAfter = filterValues
		case "maintenance­-event-not­before":
			filters.MaintenanceEventsNotBefore = filterValues
		case "vm­-state­-code":
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
