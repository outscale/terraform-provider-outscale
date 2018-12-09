package outscale

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIVMSState() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMSStateRead,
		Schema: getOAPIVMSStateDataSourceSchema(),
	}
}

func dataSourceOutscaleOAPIVMSStateRead(d *schema.ResourceData, meta interface{}) error {
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
		params.Filters.VmIds = oapiExpandStringList(instanceIds.([]interface{}))
	}

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

	if len(filteredStates) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	return statusDescriptionOAPIVMSStateAttributes(d, filteredStates)
}

func getOAPIVMSStateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"filter": dataSourceFiltersSchema(),
		"vm_id": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vm_state_set": { //events_set
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sub_region_name": { //availability_zone
						Type:     schema.TypeString,
						Computed: true,
					},
					"maintenance_event": { //events_set
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": { //code
									Type:     schema.TypeString,
									Computed: true,
								},
								"description": { //
									Type:     schema.TypeString,
									Computed: true,
								},
								"not_after": { // not_before
									Type:     schema.TypeString,
									Computed: true,
								},
								"not_before": { // not_after
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"vm_id": { //instance_id
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": { //instance_state
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"state_code": { // code
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
					"comment_item": { //details
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": { //status
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"comment_state": { //state
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
	}
}

func statusDescriptionOAPIVMSStateAttributes(d *schema.ResourceData, status []oapi.VmStates) error {

	d.SetId(resource.UniqueId())

	states := make([]map[string]interface{}, len(status))

	for k, v := range status {
		state := make(map[string]interface{})

		state["sub_region_name"] = v.SubregionName

		events := oapiEventsSet(v.MaintenanceEvents)
		state["maintenance_event"] = events

		state["state"] = v.VmState

		state["comment_item"] = v.VmState
		state["comment_state"] = v.VmState

		states[k] = state
	}

	return d.Set("vm_state_set", states)
}

func oapiExpandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
